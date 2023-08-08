package csd

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceZoneDelegation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceZoneDelegationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "FQDN of the DNS zone",
				Type:        schema.TypeString,
				Required:    true,
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
	}
}

func dataSourceZoneDelegationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	//zoneDelegation, err := apiClient.curl("GET", fmt.Sprintf("/v2/zone_delegations/%s", name), strings.NewReader(""))
	zoneDelegation, err := apiClient.getZoneDelegation(name)
	if err != nil {
		return err
	}

	nameServers := make([]interface{}, len(zoneDelegation.NameServers), len(zoneDelegation.NameServers))
	for i, item := range zoneDelegation.NameServers {
		nameServers[i] = item
	}

	// sets the response body (zone delegation object) to Terraform zone_delegation data source
	if err := d.Set("name_servers", nameServers); err != nil {
		return diag.FromErr(err)
	}

	// always run; set the resource ID to timestamp (forces this resource to refresh during every Terraform apply)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
