// Package version prints the current version of cant. Synopsis: cant version
package version

import (
	"cant/cmd"
	"cant/util/globals"
	"fmt"

	"github.com/spf13/cobra"
)

// init adds the version command as subcommand to the Root
func init() {
	cmd.Root.AddCommand(version)
}

// version command
var version = &cobra.Command{
	Use:               "version",
	Short:             "Print version",
	Long:              `Prints cant's version number`,
	PersistentPreRun:  func(cmd *cobra.Command, args []string) {}, // overwrite pre runs
	PersistentPostRun: func(cmd *cobra.Command, args []string) {}, // overwrite post runs
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		fmt.Printf("cant version %s\n\n", globals.Version)
		cmd.Root.Usage()
		return nil
	},
}
