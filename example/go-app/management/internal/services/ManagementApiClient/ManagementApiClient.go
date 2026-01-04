package ManagementApiClient

import (
	"context"
	"encoding/json"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_go_app_ManagementApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/ManagementApiClient"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	common "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/common"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/contracts/RageApiClient"
	models_api_linked_identities "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_linked_identities"
	models_api_passkey "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_passkey"
	models_api_profile "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_profile"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	fluffycore_go_app_cookies "github.com/fluffy-bunny/fluffycore/go-app/cookies"
)

type (
	service struct {
		managementClientConfigAccessor contracts_go_app_ManagementApiClient.IManagementClientConfigAccessor
		rageApiClient                  contracts_go_app_RageApiClient.IRageApiClient
		wellknownCookieNames           contracts_cookies.IWellknownCookieNames
	}
)

var stemService = (*service)(nil)

var _ contracts_go_app_ManagementApiClient.IManagementApiClient = stemService

func (s *service) Ctor(
	managementClientConfigAccessor contracts_go_app_ManagementApiClient.IManagementClientConfigAccessor,
	rageApiClient contracts_go_app_RageApiClient.IRageApiClient,
	wellknownCookieNames contracts_cookies.IWellknownCookieNames,
) (contracts_go_app_ManagementApiClient.IManagementApiClient, error) {
	return &service{
		managementClientConfigAccessor: managementClientConfigAccessor,
		rageApiClient:                  rageApiClient,
		wellknownCookieNames:           wellknownCookieNames,
	}, nil
}

func AddScopedIManagementApiClient(cb di.ContainerBuilder) {
	di.AddScoped[contracts_go_app_ManagementApiClient.IManagementApiClient](cb, stemService.Ctor)
}

// getCSRFToken retrieves the CSRF token from the _csrf cookie
func (s *service) getCSRFToken() string {

	ccCSRF := s.wellknownCookieNames.GetCookieName(contracts_cookies.CookieName_CSRF)
	csrfToken, err := fluffycore_go_app_cookies.GetCookie[string](ccCSRF)
	if err != nil {
		return ""
	}
	return csrfToken
}

// buildCustomHeaders creates custom headers map with CSRF token if it exists
func (s *service) buildCustomHeaders() map[string]string {
	csrfToken := s.getCSRFToken()
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

func (s *service) Login(ctx context.Context, request *models_api_login_models.LoginRequest) (*common.WrappedResonseT[models_api_login_models.LoginResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_login_models.LoginResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_Login),
			Data:          request,
			CustomHeaders: s.buildCustomHeaders(),
		})
}

func (s *service) Logout(ctx context.Context, request *models_api_login_models.LogoutRequest) (*common.WrappedResonseT[models_api_login_models.LogoutResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_login_models.LogoutResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_Logout),
			Data:          request,
			CustomHeaders: s.buildCustomHeaders(),
		})
}

func (s *service) PasswordResetStart(ctx context.Context, request *models_api_login_models.PasswordResetStartRequest) (*common.WrappedResonseT[models_api_login_models.PasswordResetStartResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_login_models.PasswordResetStartResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_PasswordResetStart),
			Data:          request,
			CustomHeaders: s.buildCustomHeaders(),
		})
}

func (s *service) PasswordResetFinish(ctx context.Context, request *models_api_login_models.PasswordResetFinishRequest) (*common.WrappedResonseT[models_api_login_models.PasswordResetFinishResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_login_models.PasswordResetFinishResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_PasswordResetFinish),
			Data:          request,
			CustomHeaders: s.buildCustomHeaders(),
		})
}

func (s *service) VerifyCode(ctx context.Context, request *models_api_login_models.VerifyCodeRequest) (*common.WrappedResonseT[models_api_login_models.VerifyCodeResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_login_models.VerifyCodeResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_VerifyCode),
			Data:          request,
			CustomHeaders: s.buildCustomHeaders(),
		})
}

func (s *service) GetUserProfile(ctx context.Context) (*common.WrappedResonseT[models_api_profile.Profile], error) {

	return common.HTTPFetchWrappedResponseT[models_api_profile.Profile](ctx,
		&common.CallInput{
			Method:        "GET",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_UserProfilePath),
			CustomHeaders: s.buildCustomHeaders(),
		})
}

func (s *service) UpdateUserProfile(ctx context.Context, profile *models_api_profile.Profile) (*common.WrappedResonseT[models_api_profile.Profile], error) {

	return common.HTTPFetchWrappedResponseT[models_api_profile.Profile](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_UserProfilePath),
			Data:          profile,
			CustomHeaders: s.buildCustomHeaders(),
		})
}

func (s *service) GetUserLinkedAccounts(ctx context.Context) (*common.WrappedResonseT[models_api_linked_identities.LinkedAccountsResponse], error) {

	return common.HTTPFetchWrappedResponseT[models_api_linked_identities.LinkedAccountsResponse](ctx,
		&common.CallInput{
			Method:        "GET",
			Url:           s.fixupApiPath(ctx, "/api/linked-accounts"),
			CustomHeaders: s.buildCustomHeaders(),
		})
}

func (s *service) DeleteUserLinkedAccount(ctx context.Context, identity string) (*common.WrappedResonseT[models_api_linked_identities.DeleteLinkedAccountResponse], error) {

	return common.HTTPFetchWrappedResponseT[models_api_linked_identities.DeleteLinkedAccountResponse](ctx,
		&common.CallInput{
			Method:        "DELETE",
			Url:           s.fixupApiPath(ctx, "/api/linked-accounts") + "?identity=" + identity,
			CustomHeaders: s.buildCustomHeaders(),
		})
}

// GetRageApiClient returns the core API client for direct access to TOTP and other APIs
func (s *service) GetRageApiClient() contracts_go_app_RageApiClient.IRageApiClient {
	return s.rageApiClient
}

// Typed wrapper methods for passkey operations from RageApiClient
func (s *service) RenamePasskeyHTTP(ctx context.Context, request *models_api_passkey.PasskeyRenameRequest) (*common.WrappedResonseT[models_api_passkey.PasskeyRenameResponse], error) {
	wrappedResp, err := s.rageApiClient.RenamePasskeyHTTP(ctx, request)
	if err != nil {
		return nil, err
	}

	// Convert map response to typed response
	var renameResp models_api_passkey.PasskeyRenameResponse
	if wrappedResp.Response != nil {
		data, err := json.Marshal(wrappedResp.Response)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &renameResp); err != nil {
			return nil, err
		}
	}

	return &common.WrappedResonseT[models_api_passkey.PasskeyRenameResponse]{
		Code:     wrappedResp.Code,
		Response: &renameResp,
	}, nil
}

func (s *service) GetPasskeysHTTP(ctx context.Context) (*common.WrappedResonseT[models_api_passkey.PasskeysResponse], error) {
	wrappedResp, err := s.rageApiClient.GetPasskeysHTTP(ctx)
	if err != nil {
		return nil, err
	}

	// Convert map response to typed response
	var passkeyResp models_api_passkey.PasskeysResponse
	if wrappedResp.Response != nil {
		data, err := json.Marshal(wrappedResp.Response)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &passkeyResp); err != nil {
			return nil, err
		}
	}

	return &common.WrappedResonseT[models_api_passkey.PasskeysResponse]{
		Code:     wrappedResp.Code,
		Response: &passkeyResp,
	}, nil
}

func (s *service) DeletePasskeyHTTP(ctx context.Context, request *models_api_passkey.PasskeyDeleteRequest) (*common.WrappedResonseT[models_api_passkey.PasskeyDeleteResponse], error) {
	wrappedResp, err := s.rageApiClient.DeletePasskeyHTTP(ctx, request)
	if err != nil {
		return nil, err
	}

	// Convert map response to typed response
	var deleteResp models_api_passkey.PasskeyDeleteResponse
	if wrappedResp.Response != nil {
		data, err := json.Marshal(wrappedResp.Response)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &deleteResp); err != nil {
			return nil, err
		}
	}

	return &common.WrappedResonseT[models_api_passkey.PasskeyDeleteResponse]{
		Code:     wrappedResp.Code,
		Response: &deleteResp,
	}, nil
}
