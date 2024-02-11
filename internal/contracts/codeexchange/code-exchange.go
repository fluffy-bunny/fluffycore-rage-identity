package codeexchange

import "context"

type (
	ExchangeCodeRequest struct {
		IDPSlug      string `json:"idp_slug"`
		ClientID     string `json:"client_id"`
		Code         string `json:"code"`
		CodeVerifier string `json:"code_verifier"`
		Nonce        string `json:"nonce"`
	}
	ExchangeCodeResponse struct {
		IDPSlug     string `json:"idp_slug"`
		ClientID    string `json:"client_id"`
		IdToken     string `json:"id_token"`
		AccessToken string `json:"access_token"`
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
