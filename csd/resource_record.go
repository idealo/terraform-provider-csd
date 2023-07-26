package csd

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordCreate,
		ReadContext:   resourceRecordRead,
		UpdateContext: resourceRecordUpdate,
		DeleteContext: resourceRecordDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the DNS record as FQDN",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"value": {
				Description: "Value of the DNS record (FQDN of Akamai Edgekey Hostname in case of CNAME)",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ttl": {
				Description: "Time to life for the record in seconds",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
			},
			"rrtype": {
				Description: "The type of DNS record",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	record := Record{
		Name:   d.Get("name").(string),
		RRType: d.Get("rrtype").(string),
		Value:  d.Get("value").(string),
		TTL:    d.Get("ttl").(int),
	}

	result, err := apiClient.createRecord(record)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s_%s", result.Name, result.RRType))
	if err := d.Set("name", result.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rrtype", result.RRType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("value", result.Value); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ttl", result.TTL); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()

	record, err := apiClient.getRecord(id)
	if err != nil {
		return err
	}

	// sets the response body (record object) to Terraform record data source
	if err := d.Set("name", record.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rrtype", record.RRType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("value", record.Value); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ttl", record.TTL); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Check for resource changes (only value and ttl are relevant at the moment)
	if d.HasChange("value") || d.HasChange("ttl") {
		apiClient := m.(*ApiClient)

		record := Record{
			Name:   d.Get("name").(string),
			RRType: d.Get("rrtype").(string),
			Value:  d.Get("value").(string),
			TTL:    d.Get("ttl").(int),
		}

		result, err := apiClient.updateRecord(record)
		if err != nil {
			return err
		}

		if err := d.Set("name", result.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("rrtype", result.RRType); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("value", result.Value); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("ttl", result.TTL); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceRecordRead(ctx, d, m)
}

func resourceRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()

	if err := apiClient.deleteRecord(id); err != nil {
		return err
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
