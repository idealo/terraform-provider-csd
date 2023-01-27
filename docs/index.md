---
page_title: "idealo-tools Provider"
subcategory: ""
description: |-
Terraform provider for interacting with Transport CSD product.
---

# idealo-tools Provider

The idealo-tools provider is used to interact with the CSD product.

Use the navigation to the left to read about the available resources.

## Example Usage

Do not keep your authentication password in HCL for production environments, use environment variables.

```terraform
provider "csd" {
  aws_access_key_id = "superSecret123!"
  aws_secret_access_key = "superSecret123!"
  aws_session_token = "superSecret123!"
}
```

## Schema

### Required

- `aws_access_key_id` (String, Sensitive)
- `aws_secret_access_key` (String, Sensitive)
- `aws_session_token` (String, Sensitive)
