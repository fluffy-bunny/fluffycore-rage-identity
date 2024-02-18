package email

import (
	"context"
	"html/template"
)

type (
	SendEmailRequest struct {
		ToEmail   string
		FromEmail string
		SubjectId string
		Template  string
		Data      map[string]interface{}
	}
	SendEmailResponse struct {
	}
	IEmailService interface {
		SendEmail(ctx context.Context, request *SendEmailRequest) (*SendEmailResponse, error)
	}
	RenderEmailRequest struct {
		Template string
		Data     map[string]interface{}
	}
	RenderEmailResponse struct {
		Html string
		Text string
	}
	IEmailRenderer interface {
		RenderEmail(ctx context.Context, request *RenderEmailRequest) (*RenderEmailResponse, error)
		SetTemplateEngine(templateEngine *template.Template) error
	}
)
