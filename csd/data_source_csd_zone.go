package csd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCsdZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCsdZoneRead,
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

func dataSourceCsdZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//
	name := d.Get("name").(string)

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
