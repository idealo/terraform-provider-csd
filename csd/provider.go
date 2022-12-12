package csd

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ApiEndpoint TODO: replace with production endpoint
var ApiEndpoint = "http://localhost:8080"

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"csd_zones": dataSourceCsdZones(),
		},
	}
}
