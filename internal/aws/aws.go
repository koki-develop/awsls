package aws

import (
	"context"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
)

type API struct {
	ec2Client       *ec2.Client
	resourcesClient *resourcegroupstaggingapi.Client
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

	return &API{
		resourcesClient: resourcegroupstaggingapi.NewFromConfig(awscfg),
		ec2Client:       ec2.NewFromConfig(awscfg),
	}, nil
}

type Resource struct {
	Service  string `json:"service"`
	Region   string `json:"region"`
	Resource string `json:"resource"`
	ARN      string `json:"arn"`
}

type Resources []Resource

func (rs Resources) Sort() {
	sort.SliceStable(rs, func(i, j int) bool {
		if rs[i].Service != rs[j].Service {
			return rs[i].Service < rs[j].Service
		}

		if rs[i].Region != rs[j].Region {
			return rs[i].Region < rs[j].Region
		}

		return rs[i].Resource < rs[j].Resource
	})
}

func (api *API) ListRegions() ([]string, error) {
	ipt := &ec2.DescribeRegionsInput{}
	resp, err := api.ec2Client.DescribeRegions(context.Background(), ipt)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(resp.Regions))
	for i, r := range resp.Regions {
		names[i] = *r.RegionName
	}

	return names, nil
}

func (api *API) GetResources() (Resources, error) {
	var rs Resources

	ipt := &resourcegroupstaggingapi.GetResourcesInput{}
	p := resourcegroupstaggingapi.NewGetResourcesPaginator(api.resourcesClient, ipt)
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
