/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface ApiErrorResponse {
  error?: string;
  internalCode?: string;
}

export interface ApiUserIdentityInfoUserIdentityInfo {
  email?: string;
  passkeyEligible?: boolean;
}

export interface ExternalIdpStartExternalIDPLoginRequest {
  directive: string;
  slug: string;
}

export interface ExternalIdpStartExternalIDPLoginResponse {
  redirectUri?: string;
}

export interface LoginModelsDirectiveDisplayPasswordPage {
  email?: string;
  hasPasskey?: boolean;
}

export interface LoginModelsDirectiveEmailCodeChallenge {
  code?: string;
}

export interface LoginModelsDirectiveRedirect {
  redirectUri?: string;
}

export interface LoginModelsDirectiveStartExternalLogin {
  slug?: string;
}

export interface LoginModelsLoginPasswordRequest {
  email: string;
  password: string;
}

export interface LoginModelsLoginPasswordResponse {
  directive: string;
  directiveEmailCodeChallenge?: LoginModelsDirectiveEmailCodeChallenge;
  directiveRedirect?: LoginModelsDirectiveRedirect;
  email: string;
}

export interface LoginModelsLoginPhaseOneRequest {
  email: string;
}

export interface LoginModelsLoginPhaseOneResponse {
  directive: string;
  directiveDisplayPasswordPage?: LoginModelsDirectiveDisplayPasswordPage;
  directiveEmailCodeChallenge?: LoginModelsDirectiveEmailCodeChallenge;
  directiveRedirect?: LoginModelsDirectiveRedirect;
  directiveStartExternalLogin?: LoginModelsDirectiveStartExternalLogin;
  email: string;
}

export enum LoginModelsPasswordResetErrorReason {
  PasswordResetErrorReasonNoError = 0,
  PasswordResetErrorReasonInvalidPassword = 1,
}

export interface LoginModelsPasswordResetFinishRequest {
  password: string;
  passwordConfirm: string;
}

export interface LoginModelsPasswordResetFinishResponse {
  directive: string;
  errorReason?: LoginModelsPasswordResetErrorReason;
}

export interface LoginModelsPasswordResetStartRequest {
  email: string;
}

export interface LoginModelsPasswordResetStartResponse {
  directive: string;
  directiveEmailCodeChallenge?: LoginModelsDirectiveEmailCodeChallenge;
  email: string;
}

export enum LoginModelsSignupErrorReason {
  SignupErrorReasonNoError = 0,
  SignupErrorReasonInvalidPassword = 1,
  SignupErrorReasonUserAlreadyExists = 2,
}

export interface LoginModelsSignupRequest {
  email: string;
  password: string;
}

export interface LoginModelsSignupResponse {
  directive: string;
  directiveEmailCodeChallenge?: LoginModelsDirectiveEmailCodeChallenge;
  directiveRedirect?: LoginModelsDirectiveRedirect;
  directiveStartExternalLogin?: LoginModelsDirectiveStartExternalLogin;
  email: string;
  errorReason?: LoginModelsSignupErrorReason;
  message?: string;
}

export interface LoginModelsVerifyCodeRequest {
  code: string;
}

export interface LoginModelsVerifyCodeResponse {
  directive: string;
  directiveRedirect?: LoginModelsDirectiveRedirect;
}

export interface ManifestIDP {
  slug?: string;
}

export interface ManifestManifest {
  social_idps?: ManifestIDP[];
}

export interface PasswordVerifyPasswordStrengthRequest {
  password: string;
}

export interface PasswordVerifyPasswordStrengthResponse {
  valid?: boolean;
}

export interface VerifyUsernameVerifyUsernameResponse {
  passkeyAvailable?: boolean;
  userName?: string;
}

import type {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  HeadersDefaults,
  ResponseType,
} from 'axios';
import axios from 'axios';

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams
  extends Omit<AxiosRequestConfig, 'data' | 'params' | 'url' | 'responseType'> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<
  FullRequestParams,
  'body' | 'method' | 'query' | 'path'
>;

export interface ApiConfig<SecurityDataType = unknown>
  extends Omit<AxiosRequestConfig, 'data' | 'cancelToken'> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = 'application/json',
  FormData = 'multipart/form-data',
  UrlEncoded = 'application/x-www-form-urlencoded',
  Text = 'text/plain',
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>['securityWorker'];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({
    securityWorker,
    secure,
    format,
    ...axiosConfig
  }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({
      ...axiosConfig,
      baseURL: axiosConfig.baseURL || '//localhost:9044',
    });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected mergeRequestParams(
    params1: AxiosRequestConfig,
    params2?: AxiosRequestConfig,
  ): AxiosRequestConfig {
    const method = params1.method || (params2 && params2.method);

    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...((method &&
          this.instance.defaults.headers[
            method.toLowerCase() as keyof HeadersDefaults
          ]) ||
          {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected stringifyFormItem(formItem: unknown) {
    if (typeof formItem === 'object' && formItem !== null) {
      return JSON.stringify(formItem);
    } else {
      return `${formItem}`;
    }
  }

  protected createFormData(input: Record<string, unknown>): FormData {
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      const propertyContent: any[] =
        property instanceof Array ? property : [property];

      for (const formItem of propertyContent) {
        const isFileType = formItem instanceof Blob || formItem instanceof File;
        formData.append(
          key,
          isFileType ? formItem : this.stringifyFormItem(formItem),
        );
      }

      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === 'boolean' ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = format || this.format || undefined;

    if (
      type === ContentType.FormData &&
      body &&
      body !== null &&
      typeof body === 'object'
    ) {
      body = this.createFormData(body as Record<string, unknown>);
    }

    if (
      type === ContentType.Text &&
      body &&
      body !== null &&
      typeof body !== 'string'
    ) {
      body = JSON.stringify(body);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(requestParams.headers || {}),
        ...(type && type !== ContentType.FormData
          ? { 'Content-Type': type }
          : {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title Swagger Example API
 * @version 1.0
 * @license Apache 2.0 (http://www.apache.org/licenses/LICENSE-2.0.html)
 * @termsOfService http://swagger.io/terms/
 * @baseUrl //localhost:9044
 * @contact API Support <support@swagger.io> (http://www.swagger.io/support)
 *
 * This is a sample server Petstore server.
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  wellKnown = {
    /**
     * @description get the public keys of the server.
     *
     * @tags root
     * @name JwksList
     * @summary get the public keys of the servere.
     * @request GET:/.well-known/jwks
     */
    jwksList: (params: RequestParams = {}) =>
      this.request<string, any>({
        path: `/.well-known/jwks`,
        method: 'GET',
        format: 'json',
        ...params,
      }),

    /**
     * @description get the status of server.
     *
     * @tags root
     * @name OpenidConfigurationList
     * @summary Show the status of server.
     * @request GET:/.well-known/openid-configuration
     */
    openidConfigurationList: (params: RequestParams = {}) =>
      this.request<string, any>({
        path: `/.well-known/openid-configuration`,
        method: 'GET',
        format: 'json',
        ...params,
      }),
  };
  api = {
    /**
     * @description This is the configuration of the server..
     *
     * @tags root
     * @name LoginPasswordCreate
     * @summary get the login manifest.
     * @request POST:/api/login-password
     */
    loginPasswordCreate: (
      request: LoginModelsLoginPasswordRequest,
      params: RequestParams = {},
    ) =>
      this.request<LoginModelsLoginPasswordResponse, string>({
        path: `/api/login-password`,
        method: 'POST',
        body: request,
        format: 'json',
        ...params,
      }),

    /**
     * @description This is the configuration of the server..
     *
     * @tags root
     * @name LoginPhaseOneCreate
     * @summary get the login manifest.
     * @request POST:/api/login-phase-one
     */
    loginPhaseOneCreate: (
      request: LoginModelsLoginPhaseOneRequest,
      params: RequestParams = {},
    ) =>
      this.request<LoginModelsLoginPhaseOneResponse, any>({
        path: `/api/login-phase-one`,
        method: 'POST',
        body: request,
        format: 'json',
        ...params,
      }),

    /**
     * @description This is the configuration of the server..
     *
     * @tags root
     * @name ManifestList
     * @summary get the login manifest.
     * @request GET:/api/manifest
     */
    manifestList: (params: RequestParams = {}) =>
      this.request<ManifestManifest, any>({
        path: `/api/manifest`,
        method: 'GET',
        format: 'json',
        ...params,
      }),

    /**
     * @description This is the configuration of the server..
     *
     * @tags root
     * @name PasswordResetFinishCreate
     * @summary get the login manifest.
     * @request POST:/api/password-reset-finish
     */
    passwordResetFinishCreate: (
      request: LoginModelsPasswordResetFinishRequest,
      params: RequestParams = {},
    ) =>
      this.request<LoginModelsPasswordResetFinishResponse, string>({
        path: `/api/password-reset-finish`,
        method: 'POST',
        body: request,
        format: 'json',
        ...params,
      }),

    /**
     * @description This is the configuration of the server..
     *
     * @tags root
     * @name PasswordResetStartCreate
     * @summary get the login manifest.
     * @request POST:/api/password-reset-start
     */
    passwordResetStartCreate: (
      request: LoginModelsPasswordResetStartRequest,
      params: RequestParams = {},
    ) =>
      this.request<LoginModelsPasswordResetStartResponse, any>({
        path: `/api/password-reset-start`,
        method: 'POST',
        body: request,
        format: 'json',
        ...params,
      }),

    /**
     * @description verify code
     *
     * @tags root
     * @name SignupCreate
     * @summary verify code.
     * @request POST:/api/signup
     */
    signupCreate: (
      request: LoginModelsSignupRequest,
      params: RequestParams = {},
    ) =>
      this.request<LoginModelsSignupResponse, string>({
        path: `/api/signup`,
        method: 'POST',
        body: request,
        format: 'json',
        ...params,
      }),

    /**
     * @description starts an external login ceremony with an external IDP.
     *
     * @tags root
     * @name StartExternalLoginCreate
     * @summary starts an external login ceremony with an external IDP
     * @request POST:/api/start-external-login
     */
    startExternalLoginCreate: (
      external_idp: ExternalIdpStartExternalIDPLoginRequest,
      params: RequestParams = {},
    ) =>
      this.request<ExternalIdpStartExternalIDPLoginResponse, ApiErrorResponse>({
        path: `/api/start-external-login`,
        method: 'POST',
        body: external_idp,
        format: 'json',
        ...params,
      }),

    /**
     * @description get the highlevel UserIdentityInfo post login.
     *
     * @tags root
     * @name UserIdentityInfoList
     * @summary get the highlevel UserIdentityInfo post login.
     * @request GET:/api/user-identity-info
     */
    userIdentityInfoList: (params: RequestParams = {}) =>
      this.request<ApiUserIdentityInfoUserIdentityInfo, string>({
        path: `/api/user-identity-info`,
        method: 'GET',
        format: 'json',
        ...params,
      }),

    /**
     * @description verify code
     *
     * @tags root
     * @name VerifyCodeCreate
     * @summary verify code.
     * @request POST:/api/verify-code
     */
    verifyCodeCreate: (
      request: LoginModelsVerifyCodeRequest,
      params: RequestParams = {},
    ) =>
      this.request<LoginModelsVerifyCodeResponse, string>({
        path: `/api/verify-code`,
        method: 'POST',
        body: request,
        format: 'json',
        ...params,
      }),

    /**
     * @description This is the configuration of the server..
     *
     * @tags root
     * @name VerifyPasswordStrengthCreate
     * @summary get the login manifest.
     * @request POST:/api/verify-password-strength
     */
    verifyPasswordStrengthCreate: (
      request: PasswordVerifyPasswordStrengthRequest,
      params: RequestParams = {},
    ) =>
      this.request<PasswordVerifyPasswordStrengthResponse, any>({
        path: `/api/verify-password-strength`,
        method: 'POST',
        body: request,
        format: 'json',
        ...params,
      }),

    /**
     * @description This is the configuration of the server..
     *
     * @tags root
     * @name VerifyUsernameCreate
     * @summary get the login manifest.
     * @request POST:/api/verify-username
     */
    verifyUsernameCreate: (params: RequestParams = {}) =>
      this.request<VerifyUsernameVerifyUsernameResponse, any>({
        path: `/api/verify-username`,
        method: 'POST',
        format: 'json',
        ...params,
      }),
  };
  error = {
    /**
     * @description get the error page.
     *
     * @tags root
     * @name ErrorList
     * @summary get the error page.
     * @request GET:/error
     */
    errorList: (params: RequestParams = {}) =>
      this.request<string, any>({
        path: `/error`,
        method: 'GET',
        format: 'json',
        ...params,
      }),
  };
  externalIdp = {
    /**
     * @description externalIDP.
     *
     * @tags root
     * @name ExternalIdpCreate
     * @summary todo
     * @request POST:/external-idp
     */
    externalIdpCreate: (
      query: {
        /** code */
        code: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<string, any>({
        path: `/external-idp`,
        method: 'POST',
        query: query,
        format: 'json',
        ...params,
      }),
  };
  forgotPassword = {
    /**
     * @description get the home page.
     *
     * @tags root
     * @name ForgotPasswordList
     * @summary get the home page.
     * @request GET:/forgot-password
     */
    forgotPasswordList: (
      query: {
        /** code */
        code: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<string, any>({
        path: `/forgot-password`,
        method: 'GET',
        query: query,
        format: 'json',
        ...params,
      }),

    /**
     * @description get the home page.
     *
     * @tags root
     * @name ForgotPasswordCreate
     * @summary get the home page.
     * @request POST:/forgot-password
     */
    forgotPasswordCreate: (
      query: {
        /** code */
        code: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<string, any>({
        path: `/forgot-password`,
        method: 'POST',
        query: query,
        format: 'json',
        ...params,
      }),
  };
  healthz = {
    /**
     * @description get the status of server.
     *
     * @tags root
     * @name HealthzList
     * @summary Show the status of server.
     * @request GET:/healthz
     */
    healthzList: (params: RequestParams = {}) =>
      this.request<string, any>({
        path: `/healthz`,
        method: 'GET',
        format: 'json',
        ...params,
      }),
  };
  oauth2 = {
    /**
     * @description get the home page.
     *
     * @tags root
     * @name CallbackList
     * @summary get the home page.
     * @request GET:/oauth2/callback
     */
    callbackList: (
      query: {
        /** code requested */
        code: string;
        /** state requested */
        state: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<string, any>({
        path: `/oauth2/callback`,
        method: 'GET',
        query: query,
        format: 'json',
        ...params,
      }),
  };
  oidcLogin = {
    /**
     * @description get the home page.
     *
     * @tags root
     * @name OidcLoginList
     * @summary get the home page.
     * @request GET:/oidc-login
     */
    oidcLoginList: (
      query: {
        /** code */
        code: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<string, any>({
        path: `/oidc-login`,
        method: 'GET',
        query: query,
        format: 'json',
        ...params,
      }),

    /**
     * @description get the home page.
     *
     * @tags root
     * @name OidcLoginCreate
     * @summary get the home page.
     * @request POST:/oidc-login
     */
    oidcLoginCreate: (
      query: {
        /** code */
        code: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<string, any>({
        path: `/oidc-login`,
        method: 'POST',
        query: query,
        format: 'json',
        ...params,
      }),
  };
  oidc = {
    /**
     * @description get the home page.
     *
     * @tags root
     * @name V1AuthList
     * @summary get the home page.
     * @request GET:/oidc/v1/auth
     */
    v1AuthList: (
      query: {
        /** client_id requested */
        client_id: string;
        /** response_type requested */
        response_type: string;
        /**
         * scope requested
         * @default ""openid profile email""
         */
        scope: string;
        /** state requested */
        state: string;
        /** redirect_uri requested */
        redirect_uri: string;
        /** audience requested */
        audience?: string;
        /** PKCE challenge code */
        code_challenge?: string;
        /**
         * PKCE challenge method
         * @default ""S256""
         */
        code_challenge_method?: string;
        /** acr_values requested */
        acr_values?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<string, any>({
        path: `/oidc/v1/auth`,
        method: 'GET',
        query: query,
        format: 'json',
        ...params,
      }),
  };
  token = {
    /**
     * @description OAuth2 token endpoint.
     *
     * @tags root
     * @name TokenCreate
     * @summary OAuth2 token endpoint.
     * @request POST:/token
     * @secure
     */
    tokenCreate: (
      query: {
        /** response_type requested */
        response_type: string;
        /**
         * scope requested
         * @default ""openid profile email""
         */
        scope: string;
        /** state requested */
        state: string;
        /** redirect_uri requested */
        redirect_uri: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<string, any>({
        path: `/token`,
        method: 'POST',
        query: query,
        secure: true,
        format: 'json',
        ...params,
      }),
  };
}
