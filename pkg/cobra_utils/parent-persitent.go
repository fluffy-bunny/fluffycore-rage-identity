package cobra_utils

import "github.com/spf13/cobra"

func ParentPersistentPreRunE(cmd *cobra.Command, args []string) error {
	parent := cmd.Parent()
	if parent != nil {
		if parent.PersistentPreRunE != nil {
			err := parent.PersistentPreRunE(parent, args)
			if err != nil {
				return err
			}
		} else {
			ParentPersistentPreRunE(parent, args)
		}
	}

	return nil
}
