package csd

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Zone struct {
	Name        string   `json:"name"`
	NameServers []string `json:"name_servers"`
	Owner       string   `json:"owner"`
}

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
			"aws_session_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AWS_SESSION_TOKEN", nil),
			},
		},
		//ResourcesMap: map[string]*schema.Resource{
		//	"csd_zone": resourceZone(),
		//},
		DataSourcesMap: map[string]*schema.Resource{
			"csd_zones": dataSourceZones(),
			//"csd_zone":  dataSourceZone(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	awsAccessKeyId := d.Get("aws_access_key_id").(string)
	awsSecretAccessKey := d.Get("aws_secret_access_key").(string)
	awsSessionToken := d.Get("aws_session_token").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (awsAccessKeyId != "") && (awsSecretAccessKey != "") && (awsSessionToken != "") {
		// TODO: prepare http client with signer etc
		return ApiClient{
			AuthInfo{
				AccessKeyId:     awsAccessKeyId,
				SecretAccessKey: awsSecretAccessKey,
				SessionToken:    awsSessionToken,
			},
		}, diags
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Unable to find authentication info for AWS",
		Detail:   "One of aws_access_key_id, aws_secret_access_key or aws_session_token is missing",
	})

	return nil, diags
}
