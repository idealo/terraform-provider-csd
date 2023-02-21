package csd

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support Markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string, commit string) func() *schema.Provider {
	return func() *schema.Provider {
		provider := &schema.Provider{
			// Configure terraform provider with AWS credentials
			Schema: map[string]*schema.Schema{
				"aws_access_key_id": {
					Description: "Defaults to `AWS_ACCESS_KEY_ID` environment variable",
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("AWS_ACCESS_KEY_ID", ""),
				},
				"aws_secret_access_key": {
					Description: "Defaults to `AWS_SECRET_ACCESS_KEY` environment variable",
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("AWS_SECRET_ACCESS_KEY", ""),
				},
				"aws_session_token": {
					Description: "Defaults to `AWS_SESSION_TOKEN` environment variable",
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
		}

		provider.ConfigureContextFunc = configure(version, commit, provider)

		return provider
	}
}

// configure Stores the AWS credentials from the provider configuration so later HTTP requests can use them
func configure(version string, commit string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		// Setup a User-Agent for the API client
		userAgent := p.UserAgent("terraform-provider-csd", fmt.Sprintf("%s (%s)", version, commit))

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
			UserAgent:       userAgent,
		}

		// Test the connection to find out if credentials are valid and endpoint is working
		if _, err := apiClient.getZones(); err != nil {
			diags = append(diags, err...)
		}

		return &apiClient, diags
	}
}
