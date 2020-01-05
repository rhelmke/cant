// generate webinterface
//go:generate rm -rf external/cant-ui/node_modules/
//go:generate yarn --cwd external/cant-ui install
//go:generate yarn --cwd external/cant-ui run build
//go:generate mkdir -p webserver/webinterface
//go:generate go-bindata -nocompress -nomemcopy -pkg webinterface -prefix external/cant-ui/build -o webserver/webinterface/webinterface.go external/cant-ui/build/...

// package main is the go entry point
package main

import (
	"fmt"
	// the root cmd
	"cant/cmd"
	// include all enabled commands
	"os"

	_ "cant/cmd/all"
)

// entrypoint
func main() {
	// execute cobra command chain
	if err := cmd.Root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
