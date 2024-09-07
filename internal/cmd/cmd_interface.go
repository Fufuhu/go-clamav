package cmd

import "github.com/spf13/cobra"

type CommandInterface interface {
	Run(command *cobra.Command, args []string)
}
