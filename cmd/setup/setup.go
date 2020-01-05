package cmd

// The setup command initializes cant

import (
	"cant/cmd"
	"cant/util/interactive"

	"github.com/spf13/cobra"
)

func init() {
	cmd.Root.AddCommand(setup)
}

var setup = &cobra.Command{
	Use:               "setup",
	Short:             "Interactive cant setup",
	PersistentPreRun:  func(cmd *cobra.Command, args []string) {}, // overwrite pre runs
	PersistentPostRun: func(cmd *cobra.Command, args []string) {}, // overwrite post runs
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		if err := interactive.Setup(); err != nil {
			return err
		}
		return nil
	},
}
