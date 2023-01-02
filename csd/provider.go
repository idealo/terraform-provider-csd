package csd

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ApiEndpoint TODO: replace with production endpoint
const ApiEndpoint = "http://localhost:8080"

type Zone struct {
	Name        string   `json:"name"`
	NameServers []string `json:"name_servers"`
	Owner       string   `json:"owner"`
}

type AuthInfo struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
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
		ResourcesMap: map[string]*schema.Resource{
			"csd_zone": resourceZone(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"csd_zones": dataSourceCsdZones(),
			"csd_zone":  dataSourceCsdZone(),
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
		// TODO: replace with http header generation
		return AuthInfo{
			AccessKeyId:     awsAccessKeyId,
			SecretAccessKey: awsSecretAccessKey,
			SessionToken:    awsSessionToken,
		}, diags
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Unable to find authentication info for AWS",
		Detail:   "One of aws_access_key_id, aws_secret_access_key or aws_session_token is missing",
	})

	return nil, diags
}

func providerConfigure2(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (username != "") && (password != "") {
		apiClient := ApiClient{Username: username, Password: password}

		return apiClient, diags
	}

	return nil, diags
}
