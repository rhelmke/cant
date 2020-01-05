package seed

import (
	"cant/cmd"
	"cant/util/interactive"

	"github.com/spf13/cobra"
)

var filePath string
var seedtype string

// init adds the seed command as subcommand to the Root
func init() {
	cmd.Root.AddCommand(seed)
	seed.PersistentFlags().StringVarP(&filePath, "file", "f", "spnpgn.csv", "Path to seeding file")
	seed.PersistentFlags().StringVarP(&seedtype, "type", "t", "spnpgn", "Seeding Type (spnpgn or filter)")
}

// seed command
var seed = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database",
	Long:  `This command can be used to seed the MySQL database used by cant`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		switch seedtype {
		case "spnpgn":
			if err := interactive.SeedPGNSPN(filePath); err != nil {
				return err
			}
		case "filter":
			if err := interactive.SeedFilter(); err != nil {
				return err
			}
		}
		return nil
	},
}
