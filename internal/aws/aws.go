package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
)

type API struct {
	client resourcegroupstaggingapi.GetResourcesAPIClient
}

type Config struct {
	Profile string
	Region  string
}

func New(cfg *Config) (*API, error) {
	fns := []func(*config.LoadOptions) error{}
	if cfg.Profile != "" {
		fns = append(fns, config.WithSharedConfigProfile(cfg.Profile))
	}
	if cfg.Region != "" {
		fns = append(fns, config.WithRegion(cfg.Region))
	}
	awscfg, err := config.LoadDefaultConfig(context.TODO(), fns...)
	if err != nil {
		return nil, err
	}

	svc := resourcegroupstaggingapi.NewFromConfig(awscfg)
	return &API{client: svc}, nil
}

type Resource struct {
	Service  string `json:"service"`
	Region   string `json:"region"`
	Resource string `json:"resource"`
	ARN      string `json:"arn"`
}

type Resources []Resource

func (api *API) GetResources() (Resources, error) {
	var rs Resources

	params := &resourcegroupstaggingapi.GetResourcesInput{}
	p := resourcegroupstaggingapi.NewGetResourcesPaginator(api.client, params)
	for p.HasMorePages() {
		page, err := p.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		for _, m := range page.ResourceTagMappingList {
			a, err := arn.Parse(*m.ResourceARN)
			if err != nil {
				return nil, err
			}
			rs = append(rs, Resource{
				Service:  a.Service,
				Region:   a.Region,
				Resource: a.Resource,
				ARN:      a.String(),
			})
		}
	}

	return rs, nil
}
