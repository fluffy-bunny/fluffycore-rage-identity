package email

import (
	"context"
	"net/smtp"
	"strings"

	mailyak "github.com/domodwyer/mailyak/v3"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_localizer "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/localizer"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		config            *contracts_email.EmailConfig
		localizer         contracts_localizer.ILocalizer
		emailRenderer     contracts_email.IEmailRenderer
		emailTemplateData contracts_email.IEmailTemplateData
	}
	NullEmailTemplateData struct{}
)

var stemService = (*service)(nil)
var stemNullEmailTemplateData = (*NullEmailTemplateData)(nil)

var _ contracts_email.IEmailService = stemService
var _ contracts_email.IEmailTemplateData = stemNullEmailTemplateData

func (s *NullEmailTemplateData) Ctor() (contracts_email.IEmailTemplateData, error) {
	return &NullEmailTemplateData{}, nil
}

func (s *NullEmailTemplateData) GetEmailTemplateData(request *contracts_email.GetEmailTemplateDataRequest) (*contracts_email.GetEmailTemplateDataResponse, error) {
	return &contracts_email.GetEmailTemplateDataResponse{
		Data: map[string]interface{}{},
	}, nil
}
func (s *service) Ctor(
	config *contracts_email.EmailConfig,
	emailRenderer contracts_email.IEmailRenderer,
	localizer contracts_localizer.ILocalizer,
	emailTemplateData contracts_email.IEmailTemplateData) (contracts_email.IEmailService, error) {
	return &service{
		config:            config,
		localizer:         localizer,
		emailRenderer:     emailRenderer,
		emailTemplateData: emailTemplateData,
	}, nil
}

// scoped email service due to localization
func AddScopedIEmailService(cb di.ContainerBuilder) {
	di.AddScoped[contracts_email.IEmailService](cb, stemService.Ctor)
	di.AddSingleton[contracts_email.IEmailTemplateData](cb, stemNullEmailTemplateData.Ctor)
}
func (s *service) validateSendEmailRequest(request *contracts_email.SendEmailRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if request.Data == nil {
		request.Data = make(map[string]interface{})
	}
	if fluffycore_utils.IsEmptyOrNil(request.HtmlTemplate) {
		return status.Error(codes.InvalidArgument, "HtmlTemplate is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.TextTemplate) {
		return status.Error(codes.InvalidArgument, "TextTemplate is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.SubjectId) {
		return status.Error(codes.InvalidArgument, "Subject is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.ToEmail) {
		return status.Error(codes.InvalidArgument, "ToEmail is empty")
	}
	return nil
}
func (s *service) SendEmail(ctx context.Context, request *contracts_email.SendEmailRequest) (*contracts_email.SendEmailResponse, error) {
	log := zerolog.Ctx(ctx).With().
		Str("method", "SendEmail").
		Interface("request", request).
		Logger()
	err := s.validateSendEmailRequest(request)
	if err != nil {
		return nil, err
	}
	localizer := s.localizer.GetLocalizer()
	subject, err := localizer.LocalizeMessage(&i18n.Message{ID: request.SubjectId})
	if err != nil {
		log.Error().Err(err).Msg("failed to localize subject")
		return nil, err
	}
	request.Data["LocalizeMessage"] = func(key string) string {
		message, _ := localizer.LocalizeMessage(&i18n.Message{ID: key})
		return message
	}
	finalData := make(map[string]interface{})
	getEmailTemplateDataResponse, err := s.emailTemplateData.GetEmailTemplateData(&contracts_email.GetEmailTemplateDataRequest{
		SubjectId: request.SubjectId,
	})
	if err == nil {
		finalData = getEmailTemplateDataResponse.Data
	} else {
		log.Error().Err(err).Msg("failed to get email template data")
	}
	// the original wins
	for key, value := range request.Data {
		finalData[key] = value
	}
	renderedEmail, err := s.emailRenderer.RenderEmail(ctx,
		&contracts_email.RenderEmailRequest{
			HtmlTemplate: request.HtmlTemplate,
			TextTemplate: request.TextTemplate,
			Data:         finalData,
		})
	if err != nil {
		return nil, err
	}
	if s.config.JustLogIt {
		log.Debug().Str("to", request.ToEmail).
			Str("subject", subject).Msg(renderedEmail.Text)
		return &contracts_email.SendEmailResponse{}, nil
	}
	mail := mailyak.New(
		s.config.Host,
		smtp.PlainAuth(
			s.config.Auth.PlainAuth.Identity,
			s.config.Auth.PlainAuth.Username,
			s.config.Auth.PlainAuth.Password,
			s.config.Auth.PlainAuth.Host))

	mail.To(request.ToEmail)
	mail.FromName(s.config.FromName)
	mail.From(s.config.FromEmail)
	mail.Subject(subject)
	mail.Plain().Set(string(renderedEmail.Text))
	mail.HTML().Set(string(renderedEmail.Html))
	//fmt.Println(string(renderedEmail.Html))
	err = mail.Send()
	if err != nil {
		log.Error().Err(err).Msg("failed to send email")
		return nil, err
	}

	return &contracts_email.SendEmailResponse{}, nil
}
func (s *service) validateSendSimpleEmailRequest(request *contracts_email.SendSimpleEmailRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(request.BodyId) {
		return status.Error(codes.InvalidArgument, "BodyId is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.SubjectId) {
		return status.Error(codes.InvalidArgument, "Subject is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.ToEmail) {
		return status.Error(codes.InvalidArgument, "ToEmail is empty")
	}
	return nil

}
func (s *service) SendSimpleEmail(ctx context.Context, request *contracts_email.SendSimpleEmailRequest) (*contracts_email.SendSimpleEmailResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateSendSimpleEmailRequest(request)
	if err != nil {
		return nil, err
	}
	localizer := s.localizer.GetLocalizer()
	message, err := localizer.LocalizeMessage(&i18n.Message{ID: request.BodyId})
	if err != nil {
		log.Error().Err(err).Msg("failed to localize message")
		return nil, err
	}
	for key, value := range request.Data {
		message = strings.ReplaceAll(message, "{"+key+"}", value)
	}
	_, err = s.SendEmail(ctx, &contracts_email.SendEmailRequest{
		ToEmail:      request.ToEmail,
		SubjectId:    request.SubjectId,
		HtmlTemplate: "emails/generic/index",
		TextTemplate: "emails/generic/txt",
		Data:         map[string]interface{}{"body": message},
	})
	if err != nil {
		return nil, err
	}
	return &contracts_email.SendSimpleEmailResponse{}, nil

}
