package idealo_tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strconv"
	"time"
)

func dataSourceZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceZonesRead,
		Schema: map[string]*schema.Schema{
			"zones": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name_servers": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type:     schema.TypeString,
								Computed: true,
							},
						},
						"owner": &schema.Schema{
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
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// TODO: replace with production endpoint
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/zones", "http://localhost:8080"), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	zones := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&zones)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("zones", zones); err != nil {
		return diag.FromErr(err)
	}

	// always run; set the resource ID
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
