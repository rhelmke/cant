// Package proxy is the main component of cant
package proxy

import (

	"cant/can/canbus"
	"cant/can/flowroutines/canflow"
	"cant/can/flowroutines/filter"
	"cant/can/flowroutines/isobusflow"
	"cant/cmd/run"
	"cant/util/flow"
	"cant/util/globals"
	"cant/webserver"

	// import proxy routes
	_ "cant/webserver/routing/proxy"
	//"github.com/pkg/profile"
	"github.com/spf13/cobra"
)

var mode string

// init adds the seed command as subcommand to the Root
func init() {
	run.Run.AddCommand(proxy)
	// proxy.PersistentFlags().StringVarP(&mode, "mode", "m", "filter", "operation mode, valid operations: passthru, identify, filter")
}

// run command
var proxy = &cobra.Command{
	Use:   "proxy",
	Short: "run the proxy component",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		// create flowgraph
		can0, err := canbus.NewSocket(globals.Config.Network.Interface0)
		if err != nil {
			return err
		}
		defer can0.Close()
		can1, err := canbus.NewSocket(globals.Config.Network.Interface1)
		if err != nil {
			return err
		}
		defer can1.Close()
		reader0 := canflow.CreateReader(can0)
		writer0 := canflow.CreateWriter(can0)
		reader1 := canflow.CreateReader(can1)
		writer1 := canflow.CreateWriter(can1)
		reader0.DefaultRoute(writer1)
		reader1.DefaultRoute(writer0)
		identifier0 := isobusflow.CreateIdentifier()
		identifier1 := isobusflow.CreateIdentifier()
		dataflow := flow.New()
		filter0 := filter.CreateDynamicFilter()
		filter1 := filter.CreateDynamicFilter()
		dataflow.Connect(reader0, identifier0)
		dataflow.Connect(identifier0, filter0)
		dataflow.Connect(filter0, writer1)
		dataflow.Connect(filter0, writer0)
		dataflow.Connect(reader1, identifier1)
		dataflow.Connect(identifier1, filter1)
		dataflow.Connect(filter1, writer0)
		dataflow.Connect(filter1, writer1)
		dataflow.Start()
		// run webserver
		if err := webserver.Serve(); err != nil {
			return err
		}
		dataflow.Join()
		dataflow.Stop()
		return nil
	},
}
