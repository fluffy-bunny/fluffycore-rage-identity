package EmailTemplateData

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
)

type (
	service struct {
	}
)

var stemService = (*service)(nil)
var _ contracts_email.IEmailTemplateData = stemService

func (s *service) Ctor() (contracts_email.IEmailTemplateData, error) {
	return &service{}, nil
}

func (s *service) GetEmailTemplateData(request *contracts_email.GetEmailTemplateDataRequest) (*contracts_email.GetEmailTemplateDataResponse, error) {
	return &contracts_email.GetEmailTemplateDataResponse{
		Data: map[string]interface{}{
			"account_url": "https://rage.localhost.dev",
		},
	}, nil
}

// scoped email service due to localization
func AddSingletonIEmailTemplateData(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_email.IEmailTemplateData](cb, stemService.Ctor)
}
