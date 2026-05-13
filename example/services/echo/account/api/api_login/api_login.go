package api_login

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_echo_login_handler "github.com/fluffy-bunny/fluffycore-rage-identity/example/contracts/echo/login_handler"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_selfoauth2provider "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/selfoauth2provider"
	contracts_session_with_options "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/session_with_options"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
	oauth2 "golang.org/x/oauth2"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config             *contracts_config.Config
		selfOAuth2Provider contracts_selfoauth2provider.ISelfOAuth2Provider
		session            contracts_session_with_options.ISessionWithOptions
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService
var _ contracts_echo_login_handler.ILoginHandler = stemService

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
	selfOAuth2Provider contracts_selfoauth2provider.ISelfOAuth2Provider,
	session contracts_session_with_options.ISessionWithOptions,
) (*service, error) {
	return &service{
		BaseHandler:        services_echo_handlers_base.NewBaseHandler(container, config),
		config:             config,
		selfOAuth2Provider: selfOAuth2Provider,
		session:            session,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.API_Login,
	)
	di.AddScoped[contracts_echo_login_handler.ILoginHandler](builder, stemService.Ctor)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) HandleLogin(c *echo.Context, loginRequest *login_models.LoginRequest) (*login_models.LoginResponse, error) {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	// Validate LoginRequest
	if loginRequest.ReturnURL == "" {
		log.Error().Msg("ReturnUrl is empty")
		return nil, status.Error(codes.InvalidArgument, "returnUrl is required")
	}
	s.WellknownCookies().DeleteAuthCompletedCookie(c)
	s.WellknownCookies().DeleteAuthCookie(c)

	ss, err := s.session.GetSession()
	if err != nil {
		log.Error().Err(err).Msg("s.session.GetSession")
		return nil, err
	}
	err = ss.New()
	if err != nil {
		log.Error().Err(err).Msg("ss.New")
		return nil, err
	}
	if err = ss.Save(); err != nil {
		log.Error().Err(err).Msg("ss.Save")
		return nil, err
	}
	state, err := randString(16)
	if err != nil {
		log.Error().Err(err).Msg("randString")
		return nil, err
	}
	nonce, err := randString(16)
	if err != nil {
		log.Error().Err(err).Msg("randString")
		return nil, err
	}

	// Store the LoginRequest in a cookie for the callback
	err = s.WellknownCookies().SetInsecureCookie(c,
		s.WellknownCookieNames().GetCookieName(contracts_cookies.CookieName_LoginRequest),
		&models.LoginGetRequest{
			ReturnUrl: loginRequest.ReturnURL,
		})
	if err != nil {
		log.Error().Err(err).Msg("SetInsecureCookie LoginRequest")
		return nil, err
	}

	// Store state and nonce in AccountStateCookie
	err = s.WellknownCookies().SetAccountStateCookie(c, &contracts_cookies.SetAccountStateCookieRequest{
		AccountStateCookie: &contracts_cookies.AccountStateCookie{
			State: state,
			Nonce: nonce,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetAccountStateCookie")
		return nil, err
	}

	authCodeOptions := []oauth2.AuthCodeOption{
		oidc.Nonce(nonce),
	}
	if loginRequest.AcrValues != "" {
		authCodeOptions = append(authCodeOptions, oauth2.SetAuthURLParam("acr_values", loginRequest.AcrValues))
	}

	getConfigResponse, err := s.selfOAuth2Provider.GetConfig(ctx)
	if err != nil {
		log.Error().Err(err).Msg("selfOAuth2Provider.GetConfig")
		return nil, err
	}
	oauth2Config := getConfigResponse.Config

	authRequestURL := oauth2Config.AuthCodeURL(state, authCodeOptions...)
	return &login_models.LoginResponse{
		RedirectURL: authRequestURL,
	}, nil
}

type LoginRequest struct {
	ReturnUrl string `param:"returnUrl" query:"returnUrl" form:"returnUrl" json:"returnUrl" xml:"returnUrl"`
	AcrValues string `param:"acrValues" query:"acrValues" form:"acrValues" json:"acrValues" xml:"acrValues"`
}

// API Login godoc
// @Summary Initiate OAuth2 login flow.
// @Description Initiates an OAuth2/OIDC authentication flow by generating state and nonce, and returning the authorization URL.
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "LoginRequest"
// @Success 200 {object} login_models.LoginResponse
// @Failure 400 {object} wellknown_echo.RestErrorResponse
// @Failure 500 {object} wellknown_echo.RestErrorResponse
// @Router /api/login [post]
func (s *service) Do(c *echo.Context) error {

	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	// Bind the LoginRequest
	loginRequest := &LoginRequest{}
	if err := c.Bind(loginRequest); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	loginResponse, err := s.HandleLogin(c,
		&login_models.LoginRequest{
			ReturnURL: loginRequest.ReturnUrl,
			AcrValues: loginRequest.AcrValues,
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
			default:
				return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
			}
		}
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	return c.JSONPretty(http.StatusOK, loginResponse, "  ")
}
func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
