package csd

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
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceZonesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// https://docs.aws.amazon.com/general/latest/gr/create-signed-request.html#create-canonical-request
	// https://github.com/aws/aws-sdk-go-v2/blob/main/aws/signer/v4/v4.go

	// HTTPMethod 			"GET"
	// CanonicalUri 		"/api/v1/zones"
	// CanonicalQueryString	""
	// CanonicalHeaders 	"host:csd.idealo.cloud"
	// SignedHeaders 		"host"
	// HashedPayload 		"..."

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/zones", ApiEndpoint), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	// decode the response
	zones := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&zones)
	if err != nil {
		return diag.FromErr(err)
	}

	// sets the response body (list of zone object) to Terraform zones data source
	if err := d.Set("zones", zones); err != nil {
		return diag.FromErr(err)
	}

	// always run; set the resource ID to timestamp (forces this resource to refresh during every Terraform apply)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
