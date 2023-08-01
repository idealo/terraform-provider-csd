package csd

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceZoneDelegations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceZoneDelegationsRead,
		Schema: map[string]*schema.Schema{
			"zone_delegations": {
				Description: "List of configured DNS zone delegations",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "FQDN of the DNS zone",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name_servers": {
							Description: "List of authoritative name servers for the zone",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceZoneDelegationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)
	var diags diag.Diagnostics

	results, err := apiClient.getZoneDelegations()
	if err != nil {
		return err
	}

	// convert zone struct into interface mapping
	// TODO: can we avoid this?
	zoneDelegations := make([]interface{}, len(results), len(results))
	for i, result := range results {
		zoneDelegation := make(map[string]interface{})
		zoneDelegation["name"] = result.Name
		zoneDelegation["name_servers"] = result.NameServers

		zoneDelegations[i] = zoneDelegation
	}

	if err := d.Set("zone_delegations", zoneDelegations); err != nil {
		return diag.FromErr(err)
	}

	// always run; set the resource ID to timestamp (forces this resource to refresh during every Terraform apply)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
