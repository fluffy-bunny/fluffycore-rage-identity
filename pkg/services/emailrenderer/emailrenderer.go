package emailrenderer

import (
	"context"
	"fmt"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/components"
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
var _ contracts_email.IEmailRenderer = stemService

func (s *service) Ctor(config *contracts_email.EmailConfig) (contracts_email.IEmailRenderer, error) {
	return &service{
		config: config,
	}, nil
}

func AddSingletonIEmailRenderer(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_email.IEmailRenderer](cb, stemService.Ctor)
}

func (s *service) buildEmailData(data map[string]interface{}) components.EmailData {
	emailData := components.EmailData{}
	if fn, ok := data["LocalizeMessage"].(func(string) string); ok {
		emailData.LocalizeMessage = fn
	} else {
		emailData.LocalizeMessage = func(key string) string { return key }
	}
	if url, ok := data["account_url"].(string); ok {
		emailData.AccountURL = url
	}
	if links, ok := data["headLinks"].([]components.EmailHeadLink); ok {
		emailData.HeadLinks = links
	}
	return emailData
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

	emailData := s.buildEmailData(request.Data)

	var htmlStr, textStr string

	switch request.HtmlTemplate {
	case "emails/generic/index":
		body, _ := request.Data["body"].(string)
		node := components.GenericEmailHTML(emailData, body)
		htmlStr, err = components.RenderEmailNode(node)
		if err != nil {
			return nil, err
		}
	case "emails/test/index":
		var routes []components.TestEmailRouteRow
		if errors, ok := request.Data["errors"].([]components.TestEmailRouteRow); ok {
			routes = errors
		}
		node := components.TestEmailHTML(emailData, routes)
		htmlStr, err = components.RenderEmailNode(node)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown html email template: %s", request.HtmlTemplate)
	}

	switch request.TextTemplate {
	case "emails/generic/txt":
		body, _ := request.Data["body"].(string)
		textStr = components.GenericEmailText(body)
	case "emails/test/txt":
		user, _ := request.Data["user"].(string)
		textStr = components.TestEmailText(user)
	default:
		return nil, fmt.Errorf("unknown text email template: %s", request.TextTemplate)
	}

	return &contracts_email.RenderEmailResponse{
		Html: htmlStr,
		Text: textStr,
	}, nil
}
