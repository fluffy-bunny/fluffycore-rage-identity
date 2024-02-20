package shared

import (
	"html/template"
)

type (
	Config struct {
		Port         int
		ClientId     string
		ClientSecret string
		Authority    string
		ACRValues    []string
	}
)

var (
	AppConfig    = &Config{}
	HtmlTemplate *template.Template
)

func Something(){
	HtmlTemplate.Execute(nil, nil)
}