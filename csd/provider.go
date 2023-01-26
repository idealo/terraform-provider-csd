package csd

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		// Configure terraform provider with AWS credentials
		Schema: map[string]*schema.Schema{
			"aws_access_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AWS_ACCESS_KEY_ID", ""),
			},
			"aws_secret_access_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AWS_SECRET_ACCESS_KEY", ""),
			},
			"aws_session_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AWS_SESSION_TOKEN", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"csd_zone": resourceZone(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"csd_zones": dataSourceZones(),
			"csd_zone":  dataSourceZone(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure Stores the AWS credentials from the provider configuration so later HTTP requests can use them
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	/*
		AWS credentials will be grabbed from environment variables or from the terraform provider configuration like this:
		provider "csd" {
		  aws_access_key_id     = "superSecret123!"
		  aws_secret_access_key = "superSecret123!"
		  aws_session_token     = "superSecret123!"
		}
	*/
	awsAccessKeyId := d.Get("aws_access_key_id").(string)
	awsSecretAccessKey := d.Get("aws_secret_access_key").(string)
	awsSessionToken := d.Get("aws_session_token").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// Check for the provided strings and create proper errors if they are missing
	if awsAccessKeyId == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to find authentication info for AWS",
			Detail:   "Value for aws_access_key_id is missing",
		})
	}
	if awsSecretAccessKey == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to find authentication info for AWS",
			Detail:   "Value for aws_secret_access_key is missing",
		})
	}
	if awsSessionToken == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to find authentication info for AWS",
			Detail:   "Value for aws_session_token is missing",
		})
	}

	// Create an API Client that holds the credentials and convenience function for HTTP communication
	apiClient := ApiClient{
		AccessKeyId:     awsAccessKeyId,
		SecretAccessKey: awsSecretAccessKey,
		SessionToken:    awsSessionToken,
	}

	// Test the connection to find out if credentials are valid and endpoint is working
	if _, err := apiClient.curl("GET", "/v1/zones", strings.NewReader("")); err != nil {
		diags = append(diags, err...)
	}

	return &apiClient, diags
}
