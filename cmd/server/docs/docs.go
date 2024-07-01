// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/.well-known/jwks": {
            "get": {
                "description": "get the public keys of the server.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the public keys of the servere.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/.well-known/openid-configuration": {
            "get": {
                "description": "get the status of server.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Show the status of server.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/login-password": {
            "post": {
                "description": "This is the configuration of the server..",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the login manifest.",
                "parameters": [
                    {
                        "description": "LoginPasswordRequest",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/login_models.LoginPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/login_models.LoginPasswordResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/login-phase-one": {
            "post": {
                "description": "This is the configuration of the server..",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the login manifest.",
                "parameters": [
                    {
                        "description": "LoginPhaseOneRequest",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/login_models.LoginPhaseOneRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/login_models.LoginPhaseOneResponse"
                        }
                    }
                }
            }
        },
        "/api/manifest": {
            "get": {
                "description": "This is the configuration of the server..",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the login manifest.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/manifest.Manifest"
                        }
                    }
                }
            }
        },
        "/api/password-reset-finish": {
            "post": {
                "description": "This is the configuration of the server..",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the login manifest.",
                "parameters": [
                    {
                        "description": "PasswordResetStartRequest",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/login_models.PasswordResetFinishRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/login_models.PasswordResetFinishResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/password-reset-start": {
            "post": {
                "description": "This is the configuration of the server..",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the login manifest.",
                "parameters": [
                    {
                        "description": "PasswordResetStartRequest",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/login_models.PasswordResetStartRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/login_models.PasswordResetStartResponse"
                        }
                    }
                }
            }
        },
        "/api/signup": {
            "post": {
                "description": "verify code",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "verify code.",
                "parameters": [
                    {
                        "description": "SignupRequest",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/login_models.SignupRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/login_models.SignupResponse"
                        }
                    },
                    "302": {
                        "description": "Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/start-external-login": {
            "post": {
                "description": "starts an external login ceremony with an external IDP.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "starts an external login ceremony with an external IDP",
                "parameters": [
                    {
                        "description": "StartExternalIDPLoginRequest",
                        "name": "external_idp",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/external_idp.StartExternalIDPLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/external_idp.StartExternalIDPLoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/verify-code": {
            "post": {
                "description": "verify code",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "verify code.",
                "parameters": [
                    {
                        "description": "VerifyCodeRequest",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/login_models.VerifyCodeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/login_models.VerifyCodeResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/verify-password-strength": {
            "post": {
                "description": "This is the configuration of the server..",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the login manifest.",
                "parameters": [
                    {
                        "description": "LoginPhaseOneRequest",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/password.VerifyPasswordStrengthRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/password.VerifyPasswordStrengthResponse"
                        }
                    }
                }
            }
        },
        "/api/verify-username": {
            "post": {
                "description": "This is the configuration of the server..",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the login manifest.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/verify_username.VerifyUsernameResponse"
                        }
                    }
                }
            }
        },
        "/error": {
            "get": {
                "description": "get the error page.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the error page.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/external-idp": {
            "post": {
                "description": "externalIDP.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "todo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "code",
                        "name": "code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/forgot-password": {
            "get": {
                "description": "get the home page.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the home page.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "code",
                        "name": "code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "get the home page.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the home page.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "code",
                        "name": "code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/healthz": {
            "get": {
                "description": "get the status of server.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Show the status of server.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/oauth2/callback": {
            "get": {
                "description": "get the home page.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the home page.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "code requested",
                        "name": "code",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "state requested",
                        "name": "state",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/oidc-login": {
            "get": {
                "description": "get the home page.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the home page.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "code",
                        "name": "code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "get the home page.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the home page.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "code",
                        "name": "code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/oidc/v1/auth": {
            "get": {
                "description": "get the home page.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "get the home page.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "client_id requested",
                        "name": "client_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "response_type requested",
                        "name": "response_type",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "\"openid profile email\"",
                        "description": "scope requested",
                        "name": "scope",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "state requested",
                        "name": "state",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "redirect_uri requested",
                        "name": "redirect_uri",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "audience requested",
                        "name": "audience",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "PKCE challenge code",
                        "name": "code_challenge",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "\"S256\"",
                        "description": "PKCE challenge method",
                        "name": "code_challenge_method",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "acr_values requested",
                        "name": "acr_values",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/token": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "OAuth2 token endpoint.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "OAuth2 token endpoint.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "response_type requested",
                        "name": "response_type",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "\"openid profile email\"",
                        "description": "scope requested",
                        "name": "scope",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "state requested",
                        "name": "state",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "redirect_uri requested",
                        "name": "redirect_uri",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "internalCode": {
                    "type": "string"
                }
            }
        },
        "external_idp.StartExternalIDPLoginRequest": {
            "type": "object",
            "required": [
                "directive",
                "slug"
            ],
            "properties": {
                "directive": {
                    "type": "string"
                },
                "slug": {
                    "type": "string"
                }
            }
        },
        "external_idp.StartExternalIDPLoginResponse": {
            "type": "object",
            "properties": {
                "redirectUri": {
                    "type": "string"
                }
            }
        },
        "login_models.DirectiveDisplayPasswordPage": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "hasPasskey": {
                    "type": "boolean"
                }
            }
        },
        "login_models.DirectiveEmailCodeChallenge": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                }
            }
        },
        "login_models.DirectiveRedirect": {
            "type": "object",
            "properties": {
                "redirectUri": {
                    "type": "string"
                }
            }
        },
        "login_models.DirectiveStartExternalLogin": {
            "type": "object",
            "properties": {
                "slug": {
                    "type": "string"
                }
            }
        },
        "login_models.LoginPasswordRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "login_models.LoginPasswordResponse": {
            "type": "object",
            "required": [
                "directive",
                "email"
            ],
            "properties": {
                "directive": {
                    "type": "string"
                },
                "directiveEmailCodeChallenge": {
                    "$ref": "#/definitions/login_models.DirectiveEmailCodeChallenge"
                },
                "directiveRedirect": {
                    "$ref": "#/definitions/login_models.DirectiveRedirect"
                },
                "email": {
                    "type": "string"
                }
            }
        },
        "login_models.LoginPhaseOneRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "login_models.LoginPhaseOneResponse": {
            "type": "object",
            "required": [
                "directive",
                "email"
            ],
            "properties": {
                "directive": {
                    "type": "string"
                },
                "directiveDisplayPasswordPage": {
                    "$ref": "#/definitions/login_models.DirectiveDisplayPasswordPage"
                },
                "directiveEmailCodeChallenge": {
                    "$ref": "#/definitions/login_models.DirectiveEmailCodeChallenge"
                },
                "directiveRedirect": {
                    "$ref": "#/definitions/login_models.DirectiveRedirect"
                },
                "directiveStartExternalLogin": {
                    "$ref": "#/definitions/login_models.DirectiveStartExternalLogin"
                },
                "email": {
                    "type": "string"
                }
            }
        },
        "login_models.PasswordResetErrorReason": {
            "type": "integer",
            "enum": [
                0,
                1
            ],
            "x-enum-varnames": [
                "PasswordResetErrorReason_NoError",
                "PasswordResetErrorReason_InvalidPassword"
            ]
        },
        "login_models.PasswordResetFinishRequest": {
            "type": "object",
            "required": [
                "password",
                "passwordConfirm"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "passwordConfirm": {
                    "type": "string"
                }
            }
        },
        "login_models.PasswordResetFinishResponse": {
            "type": "object",
            "required": [
                "directive"
            ],
            "properties": {
                "directive": {
                    "type": "string"
                },
                "errorReason": {
                    "$ref": "#/definitions/login_models.PasswordResetErrorReason"
                }
            }
        },
        "login_models.PasswordResetStartRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "login_models.PasswordResetStartResponse": {
            "type": "object",
            "required": [
                "directive",
                "email"
            ],
            "properties": {
                "directive": {
                    "type": "string"
                },
                "directiveEmailCodeChallenge": {
                    "$ref": "#/definitions/login_models.DirectiveEmailCodeChallenge"
                },
                "email": {
                    "type": "string"
                }
            }
        },
        "login_models.SignupErrorReason": {
            "type": "integer",
            "enum": [
                0,
                1,
                2
            ],
            "x-enum-varnames": [
                "SignupErrorReason_NoError",
                "SignupErrorReason_InvalidPassword",
                "SignupErrorReason_UserAlreadyExists"
            ]
        },
        "login_models.SignupRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "login_models.SignupResponse": {
            "type": "object",
            "required": [
                "directive",
                "email"
            ],
            "properties": {
                "directive": {
                    "type": "string"
                },
                "directiveEmailCodeChallenge": {
                    "$ref": "#/definitions/login_models.DirectiveEmailCodeChallenge"
                },
                "directiveRedirect": {
                    "$ref": "#/definitions/login_models.DirectiveRedirect"
                },
                "directiveStartExternalLogin": {
                    "$ref": "#/definitions/login_models.DirectiveStartExternalLogin"
                },
                "email": {
                    "type": "string"
                },
                "errorReason": {
                    "$ref": "#/definitions/login_models.SignupErrorReason"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "login_models.VerifyCodeRequest": {
            "type": "object",
            "required": [
                "code"
            ],
            "properties": {
                "code": {
                    "type": "string"
                }
            }
        },
        "login_models.VerifyCodeResponse": {
            "type": "object",
            "required": [
                "directive"
            ],
            "properties": {
                "directive": {
                    "type": "string"
                },
                "directiveRedirect": {
                    "$ref": "#/definitions/login_models.DirectiveRedirect"
                }
            }
        },
        "manifest.IDP": {
            "type": "object",
            "properties": {
                "slug": {
                    "type": "string"
                }
            }
        },
        "manifest.Manifest": {
            "type": "object",
            "properties": {
                "social_idps": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/manifest.IDP"
                    }
                }
            }
        },
        "password.VerifyPasswordStrengthRequest": {
            "type": "object",
            "required": [
                "password"
            ],
            "properties": {
                "password": {
                    "type": "string"
                }
            }
        },
        "password.VerifyPasswordStrengthResponse": {
            "type": "object",
            "properties": {
                "valid": {
                    "type": "boolean"
                }
            }
        },
        "verify_username.VerifyUsernameResponse": {
            "type": "object",
            "properties": {
                "passkeyAvailable": {
                    "type": "boolean"
                },
                "userName": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9044",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Swagger Example API",
	Description:      "This is a sample server Petstore server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
