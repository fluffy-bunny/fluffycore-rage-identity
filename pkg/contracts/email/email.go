package email

import (
	"context"
	"html/template"
)

type (
	SendEmailRequest struct {
		ToEmail      string
		SubjectId    string
		HtmlTemplate string
		TextTemplate string
		Data         map[string]interface{}
	}
	SendEmailResponse      struct{}
	SendSimpleEmailRequest struct {
		ToEmail   string
		SubjectId string
		BodyId    string
		Data      map[string]string
	}
	SendSimpleEmailResponse struct{}
	IEmailService           interface {
		SendEmail(ctx context.Context, request *SendEmailRequest) (*SendEmailResponse, error)
		SendSimpleEmail(ctx context.Context, request *SendSimpleEmailRequest) (*SendSimpleEmailResponse, error)
	}
	RenderEmailRequest struct {
		HtmlTemplate string
		TextTemplate string
		Data         map[string]interface{}
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
		FromName       string `json:"fromName"`
		FromEmail      string `json:"fromEmail"`
		JustLogIt      bool   `json:"justLogIt"`
	}
	IEmailRenderer interface {
		RenderEmail(ctx context.Context, request *RenderEmailRequest) (*RenderEmailResponse, error)
	}
)
