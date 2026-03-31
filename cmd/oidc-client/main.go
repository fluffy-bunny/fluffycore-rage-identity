/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	cmd "github.com/fluffy-bunny/fluffycore-rage-identity/cmd/oidc-client/cmd"
	shared "github.com/fluffy-bunny/fluffycore-rage-identity/cmd/oidc-client/shared"
)

func main() {
	shared.Version = version
	cmd.Execute()
}

var version = "Development"
