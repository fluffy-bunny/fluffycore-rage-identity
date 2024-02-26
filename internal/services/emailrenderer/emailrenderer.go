package emailrenderer

import (
	"bytes"
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/internal/contracts/email"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		config *contracts_email.EmailConfig
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_email.IEmailRenderer = stemService
}
func (s *service) Ctor(config *contracts_email.EmailConfig) (contracts_email.IEmailRenderer, error) {
	return &service{
		config: config,
	}, nil
}

func AddSingletonIEmailRenderer(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_email.IEmailRenderer](cb, stemService.Ctor)
}
func (s *service) validateRenderEmailRequest(request *contracts_email.RenderEmailRequest) error {
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
	if fluffycore_utils.IsEmptyOrNil(request.Data) {
		return status.Error(codes.InvalidArgument, "Data is empty")
	}
	return nil
}
func (s *service) RenderEmail(ctx context.Context, request *contracts_email.RenderEmailRequest) (*contracts_email.RenderEmailResponse, error) {
	err := s.validateRenderEmailRequest(request)
	if err != nil {
		return nil, err
	}
	streamWriter := new(bytes.Buffer)
	err = s.config.TemplateEngine.ExecuteTemplate(streamWriter, request.HtmlTemplate, request.Data)
	if err != nil {
		return nil, err
	}
	html := streamWriter.Bytes()

	streamWriter = new(bytes.Buffer)
	err = s.config.TemplateEngine.ExecuteTemplate(streamWriter, request.TextTemplate, request.Data)
	if err != nil {
		return nil, err
	}
	text := streamWriter.Bytes()
	return &contracts_email.RenderEmailResponse{
		Html: string(html),
		Text: string(text),
	}, nil
}
