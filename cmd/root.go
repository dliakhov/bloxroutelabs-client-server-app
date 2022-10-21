package cmd

import (
	"github.com/spf13/cobra"
)

func NewCLI() *cobra.Command {
	cli := &cobra.Command{
		Use:   "client-server",
		Short: "",
	}

	cli.AddCommand(serverCmd())
	cli.AddCommand(clientCmd())

	return cli
}
