/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package hashpassword

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	cobra_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/cobra_utils"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	cobra "github.com/spf13/cobra"
	"golang.org/x/term"
)

// hashPasswordCmd generates an Argon2id password hash suitable for, e.g., an
// operator portal OPERATOR_..._PASSWORD_HASH env var
// (see docs/OPERATOR_PORTAL_PLAN.md). This uses the exact same
// utils.GeneratePasswordHash the rest of the app uses to hash passwords, so
// the output is guaranteed compatible with contracts_identity.IPasswordHasher.
var hashPasswordCmd = &cobra.Command{
	Use:               "hash-password",
	PersistentPreRunE: cobra_utils.ParentPersistentPreRunE,
	Short:             "Generate an Argon2id password hash",
	Long: `Generate an Argon2id password hash using the same algorithm/parameters
this app uses everywhere else (utils.GeneratePasswordHash).

Pass the password as an argument, or omit it to be prompted (input hidden,
not echoed to the terminal — preferred, since a password on the command line
ends up in your shell history):

  oidc-client hash-password
  oidc-client hash-password 'my-password'
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var password string
		if len(args) > 0 {
			password = args[0]
		} else {
			var err error
			password, err = readPasswordHidden()
			if err != nil {
				return err
			}
		}
		if strings.TrimSpace(password) == "" {
			return fmt.Errorf("password must not be empty")
		}
		hash, err := utils.GeneratePasswordHash(password)
		if err != nil {
			return err
		}
		fmt.Println(hash)
		return nil
	},
}

func readPasswordHidden() (string, error) {
	fmt.Fprint(os.Stderr, "Password: ")
	fd := int(os.Stdin.Fd())
	// term.ReadPassword needs a real terminal file descriptor; fall back to
	// a plain (visible) line read if stdin isn't a terminal (e.g. piped).
	if term.IsTerminal(fd) {
		b, err := term.ReadPassword(fd)
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func InitCommand(parent *cobra.Command) {
	parent.AddCommand(hashPasswordCmd)
}
