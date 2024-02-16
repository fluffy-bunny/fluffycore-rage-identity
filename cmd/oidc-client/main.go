/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/fluffy-bunny/fluffycore-rage-oidc/cmd/oidc-client/cmd"
	"github.com/fluffy-bunny/fluffycore-rage-oidc/cmd/oidc-client/shared"
)

func main() {
	shared.Version = version
	cmd.Execute()
}

var version = "Development"
