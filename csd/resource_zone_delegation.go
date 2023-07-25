package csd

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceZoneDelegation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceZoneDelegationCreate,
		ReadContext:   resourceZoneDelegationRead,
		UpdateContext: resourceZoneDelegationUpdate,
		DeleteContext: resourceZoneDelegationDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "FQDN of the DNS zone",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name_servers": {
				Description: "List of authoritative name servers for the zone",
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

func resourceZoneDelegationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	zoneDelegation := ZoneDelegation{
		Name:        d.Get("name").(string),
		NameServers: []string{},
	}
	for _, ns := range d.Get("name_servers").([]interface{}) {
		zoneDelegation.NameServers = append(zoneDelegation.NameServers, ns.(string))
	}

	result, err := apiClient.createZoneDelegation(zoneDelegation)
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

func resourceZoneDelegationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	zoneDelegation, err := apiClient.getZoneDelegation(name)
	if err != nil {
		return err
	}

	// sets the response body (zoneDelegation object) to Terraform zone_delegation data source
	if err := d.Set("name", zoneDelegation.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name_servers", zoneDelegation.NameServers); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZoneDelegationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Check for resource changes (only name servers are relevant at the moment)
	if d.HasChange("name_servers") {
		apiClient := m.(*ApiClient)

		name := d.Id()
		zoneDelegation := ZoneDelegation{
			Name:        name,
			NameServers: []string{},
		}
		for _, ns := range d.Get("name_servers").([]interface{}) {
			zoneDelegation.NameServers = append(zoneDelegation.NameServers, ns.(string))
		}

		result, err := apiClient.updateZoneDelegation(zoneDelegation)
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

	return resourceZoneDelegationRead(ctx, d, m)
}

func resourceZoneDelegationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	if err := apiClient.deleteZoneDelegation(name); err != nil {
		return err
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
