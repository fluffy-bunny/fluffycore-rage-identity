/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	cmd "github.com/fluffy-bunny/fluffycore-rage-identity/cmd/oidc-client/cmd"
	shared "github.com/fluffy-bunny/fluffycore-rage-identity/cmd/oidc-client/shared"
	fluffycore_echo_templates "github.com/fluffy-bunny/fluffycore/echo/templates"
)

func main() {
	templateEngine, err := fluffycore_echo_templates.FindAndParseTemplates("../server/static/templates_email", nil)
	if err != nil {
		panic(err)
	}

	shared.HtmlTemplate = templateEngine
	shared.Version = version
	cmd.Execute()

}

var version = "Development"
