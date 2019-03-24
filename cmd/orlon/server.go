package main

import (
	"github.com/getpolymer/orlon/internal"
	"github.com/spf13/cobra"
)

var (
	cmdServer = &cobra.Command{
		Use:   "server",
		Short: "Start the remote server",
		Long:  "Start the remote server",
		RunE:  runCmdServer,
	}
)

func init() {
	root.AddCommand(cmdServer)
}

func runCmdServer(cmd *cobra.Command, args []string) error {
	internal.StartServer()

	return nil
}
