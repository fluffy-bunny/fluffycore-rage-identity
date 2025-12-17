package ManagementApiClient

import (
	"context"
	"encoding/json"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_go_app_ManagementApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/ManagementApiClient"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/models"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/contracts/RageApiClient"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	fluffycore_go_app_cookies "github.com/fluffy-bunny/fluffycore/go-app/cookies"
	fluffycore_go_app_fetch "github.com/fluffy-bunny/fluffycore/go-app/fetch"
)

type (
	service struct {
		managementClientConfigAccessor contracts_go_app_ManagementApiClient.IManagementClientConfigAccessor
		rageApiClient                  contracts_go_app_RageApiClient.IRageApiClient
	}
)

var stemService = (*service)(nil)

var _ contracts_go_app_ManagementApiClient.IManagementApiClient = stemService

func (s *service) Ctor(
	managementClientConfigAccessor contracts_go_app_ManagementApiClient.IManagementClientConfigAccessor,
	rageApiClient contracts_go_app_RageApiClient.IRageApiClient,
) (contracts_go_app_ManagementApiClient.IManagementApiClient, error) {
	return &service{
		managementClientConfigAccessor: managementClientConfigAccessor,
		rageApiClient:                  rageApiClient,
	}, nil
}

func AddScopedIManagementApiClient(cb di.ContainerBuilder) {
	di.AddScoped[contracts_go_app_ManagementApiClient.IManagementApiClient](cb, stemService.Ctor)
}

// getCSRFToken retrieves the CSRF token from the _csrf cookie
func getCSRFToken() string {
	csrfToken, err := fluffycore_go_app_cookies.GetCookie[string]("_csrf")
	if err != nil {
		return ""
	}
	return csrfToken
}

// buildCustomHeaders creates custom headers map with CSRF token if it exists
func buildCustomHeaders() map[string]string {
	csrfToken := getCSRFToken()
	if csrfToken != "" {
		return map[string]string{
			"X-Csrf-Token": csrfToken,
		}
	}
	return nil
}

func (s *service) fixupApiPath(ctx context.Context, relativePath string) string {
	appConfig := s.managementClientConfigAccessor.GetManagementClientConfig(ctx)
	if appConfig == nil {
		panic("AppConfig is not set in AppConfigAccessor")
	}
	return appConfig.ApiBaseUrl + relativePath
}

func (s *service) Login(ctx context.Context, request *models_api_login_models.LoginRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LoginResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.LoginResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_Login),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) Logout(ctx context.Context, request *models_api_login_models.LogoutRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LogoutResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.LogoutResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_Logout),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) PasswordResetStart(ctx context.Context, request *models_api_login_models.PasswordResetStartRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.PasswordResetStartResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.PasswordResetStartResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_PasswordResetStart),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) PasswordResetFinish(ctx context.Context, request *models_api_login_models.PasswordResetFinishRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.PasswordResetFinishResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.PasswordResetFinishResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_PasswordResetFinish),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) VerifyCode(ctx context.Context, request *models_api_login_models.VerifyCodeRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.VerifyCodeResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.VerifyCodeResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_VerifyCode),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) GetUserProfile(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models.Profile], error) {

	return fluffycore_go_app_fetch.FetchWrappedResponseT[models.Profile](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "GET",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_UserProfilePath),
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) UpdateUserProfile(ctx context.Context, profile *models.Profile) (*fluffycore_go_app_fetch.WrappedResonseT[models.Profile], error) {

	return fluffycore_go_app_fetch.FetchWrappedResponseT[models.Profile](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_UserProfilePath),
			Data:          profile,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) GetUserLinkedAccounts(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models.LinkedAccountsResponse], error) {

	return fluffycore_go_app_fetch.FetchWrappedResponseT[models.LinkedAccountsResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "GET",
			Url:           s.fixupApiPath(ctx, "/api/linked-accounts"),
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) DeleteUserLinkedAccount(ctx context.Context, identity string) (*fluffycore_go_app_fetch.WrappedResonseT[models.DeleteLinkedAccountResponse], error) {

	return fluffycore_go_app_fetch.FetchWrappedResponseT[models.DeleteLinkedAccountResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "DELETE",
			Url:           s.fixupApiPath(ctx, "/api/linked-accounts") + "?identity=" + identity,
			CustomHeaders: buildCustomHeaders(),
		})
}

// GetRageApiClient returns the core API client for direct access to TOTP and other APIs
func (s *service) GetRageApiClient() contracts_go_app_RageApiClient.IRageApiClient {
	return s.rageApiClient
}

// Typed wrapper methods for passkey operations from RageApiClient
func (s *service) RenamePasskeyHTTP(ctx context.Context, credentialID string, body *models.PasskeyRenameRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models.PasskeyRenameResponse], error) {
	wrappedResp, err := s.rageApiClient.RenamePasskeyHTTP(ctx, credentialID, body.FriendlyName)
	if err != nil {
		return nil, err
	}

	// Convert map response to typed response
	var renameResp models.PasskeyRenameResponse
	if wrappedResp.Response != nil {
		data, err := json.Marshal(wrappedResp.Response)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &renameResp); err != nil {
			return nil, err
		}
	}

	return &fluffycore_go_app_fetch.WrappedResonseT[models.PasskeyRenameResponse]{
		Code:     wrappedResp.Code,
		Response: &renameResp,
	}, nil
}

func (s *service) GetPasskeysHTTP(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models.PasskeysResponse], error) {
	wrappedResp, err := s.rageApiClient.GetPasskeysHTTP(ctx)
	if err != nil {
		return nil, err
	}

	// Convert map response to typed response
	var passkeyResp models.PasskeysResponse
	if wrappedResp.Response != nil {
		data, err := json.Marshal(wrappedResp.Response)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &passkeyResp); err != nil {
			return nil, err
		}
	}

	return &fluffycore_go_app_fetch.WrappedResonseT[models.PasskeysResponse]{
		Code:     wrappedResp.Code,
		Response: &passkeyResp,
	}, nil
}

func (s *service) DeletePasskeyHTTP(ctx context.Context, credentialID string) (*fluffycore_go_app_fetch.WrappedResonseT[models.PasskeyDeleteResponse], error) {
	wrappedResp, err := s.rageApiClient.DeletePasskeyHTTP(ctx, credentialID)
	if err != nil {
		return nil, err
	}

	// Convert map response to typed response
	var deleteResp models.PasskeyDeleteResponse
	if wrappedResp.Response != nil {
		data, err := json.Marshal(wrappedResp.Response)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &deleteResp); err != nil {
			return nil, err
		}
	}

	return &fluffycore_go_app_fetch.WrappedResonseT[models.PasskeyDeleteResponse]{
		Code:     wrappedResp.Code,
		Response: &deleteResp,
	}, nil
}
