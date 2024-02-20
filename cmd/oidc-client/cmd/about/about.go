/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package about

import (
	"fmt"

	shared "github.com/fluffy-bunny/fluffycore-rage-oidc/cmd/oidc-client/shared"
	cobra_utils "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/cobra_utils"
	cobra "github.com/spf13/cobra"
)

// aboutCmd represents the about command
var aboutCmd = &cobra.Command{
	Use:               "about",
	PersistentPreRunE: cobra_utils.ParentPersistentPreRunE,
	Short:             "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\n", shared.Version)
	},
}

func InitCommand(parent *cobra.Command) {
	parent.AddCommand(aboutCmd)
}
