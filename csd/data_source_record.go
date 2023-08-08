package csd

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceRecord() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRecordRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the DNS record as FQDN",
				Type:        schema.TypeString,
				Required:    true,
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
	}
}

func dataSourceRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)
	var diags diag.Diagnostics

	//record, err := apiClient.curl("GET", fmt.Sprintf("/v2/record/%s_%s", name, rrtype), strings.NewReader(""))
	record, err := apiClient.getRecord(d.Get("name").(string))
	if err != nil {
		return err
	}

	// sets the response body (record object) to Terraform record data source
	if err := d.Set("ttl", record.TTL); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("value", record.Value); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rrtype", record.RRType); err != nil {
		return diag.FromErr(err)
	}

	// always run; set the resource ID to timestamp (forces this resource to refresh during every Terraform apply)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
