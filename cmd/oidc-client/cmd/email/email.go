/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package email

import (
	"bytes"
	"fmt"
	"net/smtp"

	mailyak "github.com/domodwyer/mailyak/v3"
	shared "github.com/fluffy-bunny/fluffycore-rage-identity/cmd/oidc-client/shared"
	cobra_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/cobra_utils"
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

		streamWriter := new(bytes.Buffer)
		data := map[string]interface{}{
			"user":         "1234",
			"organization": "Perfect Corp",
			"account_url":  "https://accounts.google.com",
			"LocalizeMessage": func(key string) string {
				return key
			},
		}

		err := shared.HtmlTemplate.ExecuteTemplate(streamWriter, "emails/test/index", data)
		if err != nil {
			return err
		}
		bb := streamWriter.Bytes()
		fmt.Println(string(bb))
		fmt.Println(fluffycore_utils.PrettyJSON([]string{}))

		streamWriter = new(bytes.Buffer)
		err = shared.HtmlTemplate.ExecuteTemplate(streamWriter, "emails/test/txt", data)
		if err != nil {
			return err
		}
		bbPlain := streamWriter.Bytes()
		fmt.Println(string(bbPlain))

		mail := mailyak.New("localhost:25", smtp.PlainAuth("", "user", "password", "localhost"))
		mail.To("dom@itsa	llbroken.com")
		mail.From("jsmith@example.com")
		mail.FromName("Bananas for Friends")

		mail.Subject("Business proposition")

		mail.Plain().Set(string(bbPlain))
		mail.HTML().Set(string(bb))
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
