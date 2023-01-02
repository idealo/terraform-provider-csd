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
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceZonesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	zones, err := apiClient.ReadZones()
	if err != nil {
		return diag.FromErr(err)
	}

	// sets the response body (list of zone object) to Terraform zones data source
	if err := d.Set("zones", zones); err != nil {
		return diag.FromErr(err)
	}

	// always run; set the resource ID to timestamp (forces this resource to refresh during every Terraform apply)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
