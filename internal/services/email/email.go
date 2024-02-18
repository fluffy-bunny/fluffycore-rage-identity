package email

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/email"
	contracts_localizer "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/localizer"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		localizer contracts_localizer.ILocalizer
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_email.IEmailService = stemService
}
func (s *service) Ctor(localizer contracts_localizer.ILocalizer) (contracts_email.IEmailService, error) {
	return &service{
		localizer: localizer,
	}, nil
}

// scoped email service due to localization
func AddScopedIEmailService(cb di.ContainerBuilder) {
	di.AddScoped[contracts_email.IEmailService](cb, stemService.Ctor)
}
func (s *service) validateSendEmailRequest(request *contracts_email.SendEmailRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if request.Data == nil {
		request.Data = make(map[string]interface{})
	}
	if fluffycore_utils.IsEmptyOrNil(request.Template) {
		return status.Error(codes.InvalidArgument, "TemplateId is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.SubjectId) {
		return status.Error(codes.InvalidArgument, "Subject is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.ToEmail) {
		return status.Error(codes.InvalidArgument, "ToEmail is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.FromEmail) {
		return status.Error(codes.InvalidArgument, "FromEmail is empty")
	}

	return nil
}
func (s *service) SendEmail(ctx context.Context, request *contracts_email.SendEmailRequest) (*contracts_email.SendEmailResponse, error) {
	err := s.validateSendEmailRequest(request)
	if err != nil {
		return nil, err
	}
	localizer := s.localizer.GetLocalizer()

	request.Data["LocalizeMessage"] = func(key string) string {
		message, _ := localizer.LocalizeMessage(&i18n.Message{ID: key})
		return message
	}
	return &contracts_email.SendEmailResponse{}, nil
}
