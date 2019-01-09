package main

import "github.com/spf13/cobra"

var (
	root = &cobra.Command{
		Use:   "orlon [command]",
		Short: "Share/Record your terminal",
		Long:  "Share/Record your terminal",
	}
)

func main() {
	root.Execute()
}
