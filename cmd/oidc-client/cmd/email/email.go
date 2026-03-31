/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package email

import (
	"fmt"
	"net/smtp"

	mailyak "github.com/domodwyer/mailyak/v3"
	cobra_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/cobra_utils"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/components"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	cobra "github.com/spf13/cobra"
)

// aboutCmd represents the about command
var aboutCmd = &cobra.Command{
	Use:               "email",
	PersistentPreRunE: cobra_utils.ParentPersistentPreRunE,
	Short:             "email",
	Long:              `email`,
	RunE: func(cmd *cobra.Command, args []string) error {

		emailData := components.EmailData{
			AccountURL: "https://accounts.google.com",
			LocalizeMessage: func(key string) string {
				return key
			},
		}

		routes := []components.TestEmailRouteRow{} // test data
		node := components.TestEmailHTML(emailData, routes)
		htmlStr, err := components.RenderEmailNode(node)
		if err != nil {
			return err
		}
		fmt.Println(htmlStr)
		fmt.Println(fluffycore_utils.PrettyJSON([]string{}))

		plainStr := components.TestEmailText("1234")
		fmt.Println(plainStr)

		mail := mailyak.New("localhost:25", smtp.PlainAuth("", "user", "password", "localhost"))
		mail.To("dom@itsallbroken.com")
		mail.From("jsmith@example.com")
		mail.FromName("Bananas for Friends")

		mail.Subject("Business proposition")

		mail.Plain().Set(plainStr)
		mail.HTML().Set(htmlStr)
		err = mail.Send()
		if err != nil {
			return err
		}
		return nil
	},
}

func InitCommand(parent *cobra.Command) {
	parent.AddCommand(aboutCmd)
}
