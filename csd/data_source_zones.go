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
	apiClient := m.(ApiClient)

	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/zones", "https://6zrrgc0ria.execute-api.eu-central-1.amazonaws.com"), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	authorizationHeaders := signer(&apiClient.AuthInfo, request)
	request.Header.Add("X-Amz-Security-Token", apiClient.AuthInfo.SessionToken)
	request.Header.Add("X-Amz-Date", authorizationHeaders.date)
	request.Header.Add("Authorization", authorizationHeaders.authorizationHeaders)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-amz-content-sha256", fmt.Sprintf("%x", authorizationHeaders.payloadHash))

	response, err := client.Do(request)
	if err != nil {
		return diag.FromErr(err)
	}
	defer response.Body.Close()

	zones := make([]map[string]interface{}, 0)
	err = json.NewDecoder(response.Body).Decode(&zones)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("zones", zones); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
