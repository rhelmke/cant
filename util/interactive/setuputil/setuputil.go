// Package setuputil creates a wrapper for clt.InteractiveSession and config.Config for easy
// commandline configuration.
package setuputil

import (
	"cant/util/config"
	"cant/util/database"
	"cant/util/globals"
	"cant/util/networking/validation"
	"fmt"
	"strconv"

	"github.com/BTBurke/clt"
)

// SetupSession is a wrapper for clt.InteractiveSession
type SetupSession struct {
	is     *clt.InteractiveSession
	config *config.Config
}

// NewSetupSession creates a new clt.InteractiveSession-Wrapper
func NewSetupSession() *SetupSession {
	return &SetupSession{is: clt.NewInteractiveSession(), config: config.NewConfig()}
}

// Welcome the user to the interactive setup session
func (sess *SetupSession) Welcome() error {
	resp := sess.is.Say("Welcome to the CANT setup assistant. You are about to (re)configure the application.").
		Warn(clt.SStyled("This will overwrite all files in '%s'!", clt.Bold, clt.Red), globals.CantBasePath).
		AskYesNo("Are you sure you want to continue?", "n")
	if clt.IsNo(resp) {
		return fmt.Errorf("abort. shutting down")
	}
	return nil
}

// newSection indicates a new section on the commandline
func (sess *SetupSession) newSection(name string) {
	sess.is.Say(clt.SStyled("["+name+"]", clt.Bold))
}

// endSection indicates a finished section on the commandline
func (sess *SetupSession) endSection(name string) {
	sess.is.Say(clt.SStyled("Finished "+name+" configuration", clt.Bold))
}

// ConfigureMySQL asks the user for relevant MySQL connection settings
func (sess *SetupSession) ConfigureMySQL() {
	sess.newSection("MySQL Database")
	sess.config.MySQL.Host = sess.is.AskWithDefault("Host", "127.0.0.1", validation.ValidateHost)
	sess.config.MySQL.Port, _ = strconv.Atoi(sess.is.AskWithDefault("Port", "3306", validation.ValidatePort))
	sess.config.MySQL.User = sess.is.AskWithDefault("User", "cant")
	sess.config.MySQL.Password = sess.is.AskPassword()
	sess.config.MySQL.DB = sess.is.AskWithDefault("\nDatabase Name", "cant")
	sess.endSection("MySQL")
}

// ConfigureWebInterface asks the user for relevant web interface settings
func (sess *SetupSession) ConfigureWebInterface() {
	sess.newSection("Webinterface")
	sess.config.Webinterface.Host = sess.is.AskWithDefault("Bind Webinterface to", "0.0.0.0", validation.ValidateHost)
	sess.config.Webinterface.Port, _ = strconv.Atoi(sess.is.AskWithDefault("Port", "8080", validation.ValidatePort))
	sess.endSection("Webinterface")
}

// ConfigureNetwork asks the user for relevant network settings
func (sess *SetupSession) ConfigureNetwork() {
	sess.newSection("Network")
	sess.config.Network.Interface0 = sess.is.AskWithDefault("Bind CAN Listener0 Interface", "can0")
	sess.config.Network.Interface1 = sess.is.AskWithDefault("Bind CAN Listener1 Interface", "can1")
	sess.endSection("Network")
}

// AllSetupFuncs returns all functions which need to be called when configuring cant
func (sess *SetupSession) AllSetupFuncs() []func() {
	return []func(){
		sess.ConfigureMySQL,
		sess.ConfigureWebInterface,
		sess.ConfigureNetwork,
	}
}

// Save is a wrapper for cant/util/config.Config and runs additional integrity checks
func (sess *SetupSession) Save() {
	sess.newSection("Running checks")
	for _, category := range *sess.config.AllIntegrityChecks() {
		for _, check := range *category {
			spinner := clt.NewProgressSpinner("Checking " + clt.SStyled("'"+check.Name+"'", clt.Bold))
			spinner.Start()
			if err := check.Run(); err != nil {
				spinner.Fail()
				sess.is.Warn(clt.SStyled("Error: '%s'", clt.Bold, clt.Red), err)
				return
			}
			spinner.Success()
		}
	}
	sess.newSection("All checks succeeded! Saving")
	if err := sess.config.Save(globals.CantConfigPath); err != nil {
		sess.is.Warn(clt.SStyled("Error: '%s'", clt.Bold, clt.Red), err)
	}
	sess.newSection("Writing initial database scheme")
	spinner := clt.NewProgressSpinner("Creating tables")
	spinner.Start()
	db, err := database.NewMySQLConnection(sess.config.MySQL.Host, sess.config.MySQL.Port, sess.config.MySQL.User, sess.config.MySQL.Password, sess.config.MySQL.DB)
	if err != nil {
		spinner.Fail()
		sess.is.Warn(clt.SStyled("Error: '%s'", clt.Bold, clt.Red), err)
		return
	}
	defer db.Close()
	if err := database.CreateTables(db); err != nil {
		spinner.Fail()
		sess.is.Warn(clt.SStyled("Error: '%s'", clt.Bold, clt.Red), err)
		return
	}
	spinner.Success()
	sess.newSection("Installation complete! You may now run 'cant seed' for importing data to the database.")
}
