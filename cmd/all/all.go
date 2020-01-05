// Package all can be imported to attach all available commands into the commandline utility
package all

import (
	// Activated commands
	_ "cant/cmd"
	_ "cant/cmd/run"
	_ "cant/cmd/run/proxy"
	_ "cant/cmd/seed"
	_ "cant/cmd/setup"
	_ "cant/cmd/version"
)
