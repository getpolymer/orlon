package main

import (
	"fmt"

	"github.com/spf13/cobra"
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
	return fmt.Errorf("Not Implemented")
}
