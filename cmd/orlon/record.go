package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cmdRecord = &cobra.Command{
		Use:   "record",
		Short: "Record your terminal",
		Long:  "Record your terminal and save it in a file",
		Run:   runCmdRecord,
	}
)

func init() {
	root.AddCommand(cmdRecord)
}

func runCmdRecord(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("Not Implemented")
}
