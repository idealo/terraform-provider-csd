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
				Description: "FQDN of the DNS zone",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name_servers": {
				Description: "List of authoritative name servers for this zone",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	//zone, err := apiClient.curl("GET", fmt.Sprintf("/v1/zones/%s", name), strings.NewReader(""))
	zone, err := apiClient.getZone(name)
	if err != nil {
		return err
	}

	nameServers := make([]interface{}, len(zone.NameServers), len(zone.NameServers))
	for i, item := range zone.NameServers {
		nameServers[i] = item
	}

	// sets the response body (zone object) to Terraform zone data source
	if err := d.Set("name_servers", nameServers); err != nil {
		return diag.FromErr(err)
	}

	// always run; set the resource ID to timestamp (forces this resource to refresh during every Terraform apply)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
