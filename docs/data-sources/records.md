---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "csd_records Data Source - terraform-provider-csd"
subcategory: ""
description: |-
  
---

# csd_records (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `records` (List of Object) List of configured DNS records (see [below for nested schema](#nestedatt--records))

<a id="nestedatt--records"></a>
### Nested Schema for `records`

Read-Only:

- `name` (String)
- `rrtype` (String)
- `ttl` (Number)
- `value` (String)
