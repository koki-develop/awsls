package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "awsls",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return err
		}

		svc := resourcegroupstaggingapi.NewFromConfig(cfg)

		var as []arn.ARN

		params := &resourcegroupstaggingapi.GetResourcesInput{}
		p := resourcegroupstaggingapi.NewGetResourcesPaginator(svc, params)
		for p.HasMorePages() {
			page, err := p.NextPage(context.Background())
			if err != nil {
				return err
			}
			for _, m := range page.ResourceTagMappingList {
				a, err := arn.Parse(*m.ResourceARN)
				if err != nil {
					return err
				}
				as = append(as, a)
			}
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(as); err != nil {
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
