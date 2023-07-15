package cmd

import (
	"os"

	"github.com/koki-develop/awsls/internal/aws"
	"github.com/koki-develop/awsls/internal/printer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "awsls",
	RunE: func(cmd *cobra.Command, args []string) error {
		api, err := aws.New(&aws.Config{})
		if err != nil {
			return err
		}

		p, err := printer.New("json")
		if err != nil {
			return err
		}

		rs, err := api.GetResources()
		if err != nil {
			return err
		}

		if err := p.Print(os.Stdout, rs); err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
