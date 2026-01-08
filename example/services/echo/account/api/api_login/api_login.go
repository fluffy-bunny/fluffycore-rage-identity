package api_login

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	shared "github.com/fluffy-bunny/fluffycore-rage-identity/cmd/oidc-client/shared"
	contracts_echo_login_handler "github.com/fluffy-bunny/fluffycore-rage-identity/example/contracts/echo/login_handler"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_session_with_options "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/session_with_options"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	oauth2 "golang.org/x/oauth2"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config *contracts_config.Config

		session contracts_session_with_options.ISessionWithOptions
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService
var _ contracts_echo_login_handler.ILoginHandler = stemService

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
	session contracts_session_with_options.ISessionWithOptions,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
		config:      config,
		session:     session,
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

var (
	callbackPath = wellknown_echo.AccountCallbackPath
)

func (s *service) HandleLogin(c echo.Context, loginRequest *login_models.LoginRequest) (*login_models.LoginResponse, error) {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	r := c.Request()

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
	if len(shared.AppConfig.ACRValues) > 0 {
		authCodeOptions = append(authCodeOptions, AcrValues(shared.AppConfig.ACRValues...))
	}

	provider, err := oidc.NewProvider(ctx, s.config.OIDCConfig.BaseUrl)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query provider.")
		return nil, err
	}

	// Build redirect URL from the request itself
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	redirectUrl := fmt.Sprintf("%s://%s%s", scheme, r.Host, callbackPath)

	config := oauth2.Config{
		ClientID:     shared.AppConfig.ClientId,
		ClientSecret: shared.AppConfig.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectUrl,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", oidc.ScopeOfflineAccess},
	}

	authRequestURL := config.AuthCodeURL(state, authCodeOptions...)
	return &login_models.LoginResponse{
		RedirectURL: authRequestURL,
	}, nil
}

type LoginRequest struct {
	ReturnUrl string `param:"returnUrl" query:"returnUrl" form:"returnUrl" json:"returnUrl" xml:"returnUrl"`
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
func (s *service) Do(c echo.Context) error {

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

func AcrValues(acr ...string) oauth2.AuthCodeOption {
	acrValues := strings.Join(acr, " ")
	return oauth2.SetAuthURLParam("acr_values", acrValues)
}
