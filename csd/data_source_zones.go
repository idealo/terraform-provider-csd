package csd

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceZonesRead,
		Schema: map[string]*schema.Schema{
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name_servers": {
							Type:     schema.TypeList,
							Computed: true,
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

func dataSourceZonesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)
	var diags diag.Diagnostics

	results, err := apiClient.getZones()
	if err != nil {
		return err
	}

	// convert zone struct into interface mapping
	// TODO: can we avoid this?
	zones := make([]interface{}, len(results), len(results))
	for i, result := range results {
		zone := make(map[string]interface{})
		zone["name"] = result.Name
		zone["name_servers"] = result.NameServers

		zones[i] = zone
	}

	if err := d.Set("zones", zones); err != nil {
		return diag.FromErr(err)
	}

	// always run; set the resource ID to timestamp (forces this resource to refresh during every Terraform apply)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
