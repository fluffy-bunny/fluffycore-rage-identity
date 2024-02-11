package codeexchange

import "context"

type (
	ExchangeCodeResponse struct {
		IdToken     string `json:"id_token"`
		AccessToken string `json:"access_token"`
	}
	ExchangeCodeRequest struct {
		Code         string `json:"code"`
		CodeVerifier string `json:"code_verifier"`
	}
	ICodeExchange interface {
		ExchangeCode(ctx context.Context, req *ExchangeCodeRequest) (*ExchangeCodeResponse, error)
	}
	IGithubCodeExchange interface {
		ICodeExchange
	}
	IGenericOIDCCodeExchange interface {
		ICodeExchange
	}
)
