---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "csd_records Data Source - terraform-provider-csd"
subcategory: ""
description: |-
  
---

# csd_records (Data Source)



## Example Usage

```terraform
data "csd_records" "all" {}

output "all_records" {
  value = data.csd_records.all
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `records` (List of Object) List of configured DNS records (see [below for nested schema](#nestedatt--records))

<a id="nestedatt--records"></a>
### Nested Schema for `records`

Read-Only:

- `id` (String) The ID of this resource.
- `name` (String) Name of the DNS record as FQDN.
- `value` (String) Value of the DNS record (FQDN of Akamai Edgekey Hostname in case of CNAME).
- `ttl` (Int) Time to life for the record in seconds.
- `rrtype` (String) The type of DNS record.