package csd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
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
			"owner": {
				Type:     schema.TypeString,
				Required: true,
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
	apiClient := m.(*ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	zone := Zone{
		Name:        d.Get("name").(string),
		NameServers: []string{},
		Owner:       d.Get("owner").(string),
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
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/zones/%s", ApiEndpoint, name), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	// decode the response
	zone := Zone{}
	err = json.NewDecoder(r.Body).Decode(&zone)
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
		// TODO: get authentication infos for http connection
		// clientThing := m.(*hc.Client)

		client := &http.Client{Timeout: 10 * time.Second}

		name := d.Id()
		zone := Zone{
			Name:        name,
			NameServers: []string{},
			Owner:       d.Get("owner").(string),
		}
		for _, ns := range d.Get("name_servers").([]interface{}) {
			zone.NameServers = append(zone.NameServers, ns.(string))
		}

		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(zone)
		if err != nil {
			return diag.FromErr(err)
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/zones/%s", ApiEndpoint, name), buf)
		if err != nil {
			return diag.FromErr(err)
		}

		r, err := client.Do(req)
		if err != nil {
			return diag.FromErr(err)
		}
		defer r.Body.Close()

		// TODO: remove if unnecessary
		if err := d.Set("last_updated", time.Now().Format(time.RFC850)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceZoneRead(ctx, d, m)
}

func resourceZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/zones/%s", ApiEndpoint, name), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
