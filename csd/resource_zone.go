package csd

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceZoneCreate,
		ReadContext:   resourceZoneRead,
		UpdateContext: resourceZoneUpdate,
		DeleteContext: resourceZoneDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "FQDN of the DNS zone",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name_servers": {
				Description: "List of authoritative name servers for this zone",
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    2,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	zone := Zone{
		Name:        d.Get("name").(string),
		NameServers: []string{},
	}
	for _, ns := range d.Get("name_servers").([]interface{}) {
		zone.NameServers = append(zone.NameServers, ns.(string))
	}

	result, err := apiClient.createZone(zone)
	if err != nil {
		return err
	}

	d.SetId(result.Name)
	if err := d.Set("name", result.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name_servers", result.NameServers); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	zone, err := apiClient.getZone(name)
	if err != nil {
		return err
	}

	// sets the response body (zone object) to Terraform zone data source
	if err := d.Set("name", zone.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name_servers", zone.NameServers); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Check for resource changes (only name servers are relevant at the moment)
	if d.HasChange("name_servers") {
		apiClient := m.(*ApiClient)

		name := d.Id()
		zone := Zone{
			Name:        name,
			NameServers: []string{},
		}
		for _, ns := range d.Get("name_servers").([]interface{}) {
			zone.NameServers = append(zone.NameServers, ns.(string))
		}

		result, err := apiClient.updateZone(zone)
		if err != nil {
			return err
		}

		if err := d.Set("name", result.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("name_servers", result.NameServers); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceZoneRead(ctx, d, m)
}

func resourceZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	if err := apiClient.deleteZone(name); err != nil {
		return err
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
