/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package version

import (
	"fmt"

	shared "github.com/fluffy-bunny/fluffycore-rage-oidc/cmd/oidc-client/shared"
	cobra_utils "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/cobra_utils"
	cobra "github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:               "version",
	Short:             "version of the cli",
	Long:              ``,
	PersistentPreRunE: cobra_utils.ParentPersistentPreRunE,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(shared.Version)
	},
}

func InitCommand(parent *cobra.Command) {
	parent.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
