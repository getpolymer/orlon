package main

import (
	"github.com/spf13/cobra"

	"github.com/getpolymer/orlon/internal"
)

var (
	cmdShare = &cobra.Command{
		Use:   "share",
		Short: "Share your terminal",
		Long:  "Share your terminal to the web",
		RunE:  runCmdShare,
	}
)

func init() {
	root.AddCommand(cmdShare)
}

func runCmdShare(cmd *cobra.Command, args []string) error {
	internal.RunPseudoTerminal()

	return nil
}
