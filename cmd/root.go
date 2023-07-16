package cmd

import (
	"os"

	"github.com/koki-develop/awsls/internal/aws"
	"github.com/koki-develop/awsls/internal/printer"
	"github.com/spf13/cobra"
)

var (
	flagProfile    string
	flagRegion     []string
	flagAllRegions bool
	flagFormat     string
)

var rootCmd = &cobra.Command{
	Use: "awsls",
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := printer.New(flagFormat)
		if err != nil {
			return err
		}

		if flagAllRegions {
			cfg := &aws.Config{Profile: flagProfile}
			api, err := aws.New(cfg)
			if err != nil {
				return err
			}
			rs, err := api.ListRegions()
			if err != nil {
				return err
			}
			flagRegion = rs
		} else {
			if len(flagRegion) == 0 {
				flagRegion = []string{""}
			}
		}

		rsrcs := aws.Resources{}
		for _, r := range flagRegion {
			cfg := &aws.Config{
				Profile: flagProfile,
				Region:  r,
			}
			api, err := aws.New(cfg)
			if err != nil {
				return err
			}

			rs, err := api.GetResources()
			if err != nil {
				return err
			}

			rsrcs = append(rsrcs, rs...)
		}

		rsrcs.Sort()
		if err := p.Print(os.Stdout, rsrcs); err != nil {
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

func init() {
	rootCmd.Flags().StringVarP(&flagProfile, "profile", "p", "", "AWS profile")

	rootCmd.Flags().StringSliceVarP(&flagRegion, "region", "r", []string{}, "AWS region")
	rootCmd.Flags().BoolVar(&flagAllRegions, "all-regions", false, "Get resources in all regions")
	rootCmd.MarkFlagsMutuallyExclusive("region", "all-regions")

	rootCmd.Flags().StringVarP(&flagFormat, "format", "f", "table", "Output format (table|json|yaml|csv)")
}
