package run

import (
	"cant/cmd"

	"github.com/spf13/cobra"
)

// init adds the seed command as subcommand to the Root
func init() {
	cmd.Root.AddCommand(Run)
}

// Run command
var Run = &cobra.Command{
	Use:   "run",
	Short: "run the main components",
}
