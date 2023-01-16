package csd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
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

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(zone); err != nil {
		return diag.FromErr(err)
	}
	if _, err := apiClient.curl("POST", "/v1/zones", buffer); err != nil {
		return err
	}

	d.SetId(zone.Name)

	resourceZoneRead(ctx, d, m)

	return diags
}

func resourceZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	zone, err := apiClient.curl("GET", fmt.Sprintf("/v1/zones/%s", name), strings.NewReader(""))
	if err != nil {
		return err
	}

	// sets the response body (zone object) to Terraform zone data source
	if err := d.Set("name", zone.(Zone).Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name_servers", zone.(Zone).NameServers); err != nil {
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

		buffer := new(bytes.Buffer)
		if err := json.NewEncoder(buffer).Encode(zone); err != nil {
			return diag.FromErr(err)
		}
		if _, err := apiClient.curl("PUT", fmt.Sprintf("/v1/zones/%s", name), buffer); err != nil {
			return err
		}
	}

	return resourceZoneRead(ctx, d, m)
}

func resourceZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	if _, err := apiClient.curl("DELETE", fmt.Sprintf("/v1/zones/%s", name), strings.NewReader("")); err != nil {
		return err
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
