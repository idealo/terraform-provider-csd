package csd

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRecordsRead,
		Schema: map[string]*schema.Schema{
			"records": {
				Description: "List of configured DNS records",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Name of the DNS record as FQDN",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"value": {
							Description: "Value of the DNS record (FQDN of Akamai Edgekey Hostname in case of CNAME)",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"ttl": {
							Description: "Time to life for the record in seconds",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"rrtype": {
							Description: "The type of DNS record",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRecordsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)
	var diags diag.Diagnostics

	results, err := apiClient.getRecords()
	if err != nil {
		return err
	}

	// convert record struct into interface mapping
	// TODO: can we avoid this?
	records := make([]interface{}, len(results), len(results))
	for i, result := range results {
		record := make(map[string]interface{})
		record["name"] = result.Name
		record["rrtype"] = result.RRType
		record["value"] = result.Value
		record["ttl"] = result.TTL

		records[i] = record
	}

	if err := d.Set("records", records); err != nil {
		return diag.FromErr(err)
	}

	// always run; set the resource ID to timestamp (forces this resource to refresh during every Terraform apply)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
