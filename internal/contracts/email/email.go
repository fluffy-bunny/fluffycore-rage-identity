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

	PlainAuth struct {
		Identity string `json:"identity"`
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
	}
	Auth struct {
		PlainAuth *PlainAuth `json:"plainAuth"`
	}
	EmailConfig struct {
		TemplateEngine *template.Template
		Host           string `json:"host"`
		Auth           *Auth  `json:"auth"`
	}
	IEmailRenderer interface {
		RenderEmail(ctx context.Context, request *RenderEmailRequest) (*RenderEmailResponse, error)
	}
)
