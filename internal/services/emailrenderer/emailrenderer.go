package emailrenderer

import (
	"context"
	"html/template"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/email"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		templateEngine *template.Template
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_email.IEmailRenderer = stemService
}
func (s *service) Ctor() (contracts_email.IEmailRenderer, error) {
	return &service{}, nil
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
	if fluffycore_utils.IsEmptyOrNil(request.Template) {
		return status.Error(codes.InvalidArgument, "TemplateId is empty")
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
	return &contracts_email.RenderEmailResponse{}, nil
}
func (s *service) SetTemplateEngine(templateEngine *template.Template) error {
	if templateEngine == nil {
		return status.Error(codes.InvalidArgument, "templateEngine is nil")
	}
	s.templateEngine = templateEngine
	return nil
}
