package csd

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ApiEndpoint TODO: replace with production endpoint
var ApiEndpoint = "http://localhost:8080"

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"aws_access_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AWS_ACCESS_KEY_ID", nil),
			},
			"aws_secret_access_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AWS_SECRET_ACCESS_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"csd_zones": dataSourceCsdZones(),
			"csd_zone":  dataSourceCsdZone(),
		},
		//ConfigureContextFunc: providerConfigure,
	}
}

//func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
//	awsAccessKeyId := d.Get("aws_access_key_id").(string)
//	awsSecretAccessKey := d.Get("aws_secret_access_key").(string)
//
//	// Warning or errors can be collected in a slice type
//	var diags diag.Diagnostics
//
//	if (awsAccessKeyId != "") && (awsSecretAccessKey != "") {
//		// TODO: replace with http header generation
//		c, err := hashicups.NewClient(nil, &awsAccessKeyId, &awsSecretAccessKey)
//		if err != nil {
//			return nil, diag.FromErr(err)
//		}
//
//		return c, diags
//	}
//
//	c, err := hashicups.NewClient(nil, nil, nil)
//	if err != nil {
//		return nil, diag.FromErr(err)
//	}
//
//	return c, diags
//}
