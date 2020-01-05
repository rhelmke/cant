// Package interactive implements various wrappers for clt.InteractiveSession for convenient commandline-dialogues
package interactive

import (
	"cant/util/interactive/setuputil"
)

// Setup coordinates the setup routine using setuputil
func Setup() error {
	// create a new interactive session
	sess := setuputil.NewSetupSession()
	// welcome the user
	err := sess.Welcome()
	if err != nil {
		return err
	}
	// run all setup functions
	for _, fn := range sess.AllSetupFuncs() {
		fn()
	}
	// save the session to a config file
	sess.Save()
	return nil
}
