// Package cmd implements a commandline utility based on https://github.com/spf13/cobra.
// This is the entrypoint for cant. Every provided functionality will be invoked by a subcommand
// of this package.
//
// cant uses viper to parse json configuration files.
package cmd

import (
	"fmt"
	"os"
	"path"

	"cant/models/filter"
	"cant/models/pgn"
	"cant/models/spn"
	"cant/util/config"
	"cant/util/database"
	"cant/util/globals"

	"github.com/BTBurke/clt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Root will be executed when no subcommand has been given
var Root = &cobra.Command{
	Use:   "cant",
	Short: "privacy-enhancing proxy for CAN-based networks",
	Long: `cant is a proxy for ISOBUS networks and is useful to protect your privacy by 
manipulating packets containing sensitive information before they reach their target.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		readConfig()
		var err error
		globals.DB, err = database.NewMySQLConnection(globals.Config.MySQL.Host, globals.Config.MySQL.Port, globals.Config.MySQL.User, globals.Config.MySQL.Password, globals.Config.MySQL.DB)
		if err != nil {
			ExitWithError(err)
		}
		pgn.Prepare()
		spn.Prepare()
		filter.Prepare()
		if err = pgn.BuildCache(); err != nil {
			ExitWithError(err)
		}
		if err = filter.BuildCache(); err != nil {
			ExitWithError(err)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if globals.DB != nil {
			globals.DB.Close()
		}
	},
}

// init will be called by the runtime when initializing
func init() {
	initPaths()
	initConfig()
}

// initPaths sets global file-/directory-path variables
func initPaths() {
	// get home directory
	var err error
	globals.UserHomePath, err = homedir.Dir()
	if err != nil {
		ExitWithError(err)
	}
	// set all paths
	globals.CantBasePath = path.Join(globals.UserHomePath, ".config", "cant")
	globals.CantConfigPath = path.Join(globals.CantBasePath, "config.json")
}

// initConfig configures viper
func initConfig() {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(globals.CantBasePath)
}

// checkPaths checks whether all global paths are existing and readable
func checkPaths() {
	// iterate over all paths and 'stat' it
	paths := []string{globals.UserHomePath, globals.CantBasePath, globals.CantConfigPath}
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			ExitWithError(fmt.Errorf("Did you run 'cant setup'? (%s)", err))
		}
	}
}

// readConfig instructs viper to read the config file
func readConfig() {
	if err := viper.ReadInConfig(); err != nil {
		ExitWithError(fmt.Errorf("Could not read config file (%s)", err))
	}
	globals.Config = config.NewConfig()
	if err := viper.Unmarshal(globals.Config); err != nil {
		ExitWithError(fmt.Errorf("Unmarshalling config failed (%s)", err))
	}
}

// ExitWithError stops the runtime with a red-colored error message
func ExitWithError(err error) {
	fmt.Fprintf(os.Stderr, clt.SStyled("Error: %v\n", clt.Red, clt.Bold), err)
	os.Exit(1)
}

// ExitWithWarning stops the runtime with a warning message, but wont return a value != 0
func ExitWithWarning(err error) {
	fmt.Printf(clt.SStyled("%s\n", clt.Bold), err)
	os.Exit(0)
}
