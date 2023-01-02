package csd

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceZoneCreate,
		ReadContext:   resourceZoneRead,
		UpdateContext: resourceZoneUpdate,
		DeleteContext: resourceZoneDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name_servers": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 2,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	zone := Zone{
		Name:        d.Get("name").(string),
		NameServers: []string{},
	}
	for _, ns := range d.Get("name_servers").([]interface{}) {
		zone.NameServers = append(zone.NameServers, ns.(string))
	}

	if err := apiClient.CreateZone(zone); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(zone.Name)

	resourceZoneRead(ctx, d, m)

	return diags
}

func resourceZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	zone, err := apiClient.ReadZone(name)
	if err != nil {
		return diag.FromErr(err)
	}

	// sets the response body (zone object) to Terraform zone data source
	if err := d.Set("name", zone.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name_servers", zone.NameServers); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner", zone.Owner); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Check for resource changes (only name servers are relevant at the moment)
	if d.HasChange("name_servers") {
		apiClient := m.(ApiClient)

		name := d.Id()
		zone := Zone{
			Name:        name,
			NameServers: []string{},
		}
		for _, ns := range d.Get("name_servers").([]interface{}) {
			zone.NameServers = append(zone.NameServers, ns.(string))
		}

		if err := apiClient.UpdateZone(zone); err != nil {
			return diag.FromErr(err)
		}

		// TODO: remove if unnecessary
		if err := d.Set("last_updated", time.Now().Format(time.RFC850)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceZoneRead(ctx, d, m)
}

func resourceZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	if err := apiClient.DeleteZone(name); err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
