package api_login

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	shared "github.com/fluffy-bunny/fluffycore-rage-identity/cmd/oidc-client/shared"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_session_with_options "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/session_with_options"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
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
		config           *contracts_config.Config
		wellknownCookies contracts_cookies.IWellknownCookies
		session          contracts_session_with_options.ISessionWithOptions
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	session contracts_session_with_options.ISessionWithOptions,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		config:           config,
		wellknownCookies: wellknownCookies,
		session:          session,
	}, nil
}

// 	API_Logout             = "/api/logout"

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
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

// API Login godoc
// @Summary Initiate OAuth2 login flow.
// @Description Initiates an OAuth2/OIDC authentication flow by generating state and nonce, and returning the authorization URL.
// @Tags authentication
// @Accept json
// @Produce json
// @Success 200 {object} login_models.LoginResponse
// @Failure 500 {object} wellknown_echo.RestErrorResponse
// @Router /api/login [get]
func (s *service) Do(c echo.Context) error {

	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	w := c.Response().Writer
	r := c.Request()
	s.wellknownCookies.DeleteAuthCookie(c)

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

	setCallbackCookie(w, r, "state", state)
	setCallbackCookie(w, r, "nonce", nonce)
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
func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: false,
	}
	http.SetCookie(w, c)
}
func AcrValues(acr ...string) oauth2.AuthCodeOption {
	acrValues := strings.Join(acr, " ")
	return oauth2.SetAuthURLParam("acr_values", acrValues)
}
