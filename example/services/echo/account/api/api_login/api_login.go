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
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_session_with_options "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/session_with_options"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	oauth2 "golang.org/x/oauth2"
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

// 	API_Logout             = "/api/logout"

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.API_Login,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

var (
	callbackPath = wellknown_echo.AccountCallbackPath
)

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

	r := c.Request()

	// Bind the LoginRequest
	loginRequest := &LoginRequest{}
	if err := c.Bind(loginRequest); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	// Validate LoginRequest
	if loginRequest.ReturnUrl == "" {
		log.Error().Msg("ReturnUrl is empty")
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: "returnUrl is required"}, "  ")
	}

	s.WellknownCookies().DeleteAuthCompletedCookie(c)
	s.WellknownCookies().DeleteAuthCookie(c)
	s.WellknownCookies().DeleteSSOCookie(c)

	ss, err := s.session.GetSession()
	if err != nil {
		log.Error().Err(err).Msg("s.session.GetSession")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	err = ss.New()
	if err != nil {
		log.Error().Err(err).Msg("ss.New")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	state, err := randString(16)
	if err != nil {
		log.Error().Err(err).Msg("randString")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	nonce, err := randString(16)
	if err != nil {
		log.Error().Err(err).Msg("randString")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	// Store the LoginRequest in a cookie for the callback
	err = s.WellknownCookies().SetInsecureCookie(c,
		s.WellknownCookieNames().GetCookieName(contracts_cookies.CookieName_LoginRequest),
		&models.LoginGetRequest{
			ReturnUrl: loginRequest.ReturnUrl,
		})
	if err != nil {
		log.Error().Err(err).Msg("SetInsecureCookie LoginRequest")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
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
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
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
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
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

	return c.JSONPretty(http.StatusOK, &login_models.LoginResponse{
		RedirectURL: authRequestURL,
	}, "  ")
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
