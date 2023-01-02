package csd

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceZoneRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
	}
}

func dataSourceZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//
	name := d.Get("name").(string)

	zone, err := apiClient.ReadZone(name)
	if err != nil {
		return diag.FromErr(err)
	}

	// sets the response body (zone object) to Terraform zone data source
	if err := d.Set("name_servers", zone.NameServers); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner", zone.Owner); err != nil {
		return diag.FromErr(err)
	}

	// always run; set the resource ID to timestamp (forces this resource to refresh during every Terraform apply)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
