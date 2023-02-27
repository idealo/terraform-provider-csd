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
				"profile": {
					Type:     schema.TypeString,
					Optional: true,
					Description: "The profile for API operations. If not set, the default profile\n" +
						"created with `aws configure` will be used.",
				},
				"region": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "eu-central-1",
					Description: "The region where AWS operations will take place. Examples\n" +
						"are us-east-1, us-west-2, etc.",
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

		// AWS credentials will be grabbed from environment variables or from ~/.aws/credentials

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		creds, err := getCreds(d.Get("profile").(string), d.Get("region").(string))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to find authentication info for AWS",
				Detail:   "Please configure your AWS credentials at ~/.aws/credentials or as enviromental variables.",
			})
		}

		awsAccessKeyId := creds.AccessKeyID
		awsSecretAccessKey := creds.SecretAccessKey
		awsSessionToken := creds.SessionToken

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
