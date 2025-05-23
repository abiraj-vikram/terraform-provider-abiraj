---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "securden_add_account Data Source - terraform-provider-securden"
subcategory: ""
description: |-
  Defines the structure for managing accounts in Securden
---

# securden_add_account (Data Source)

Defines the structure for managing accounts in Securden



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_name` (String) The name associated with the account
- `account_title` (String) The title associated with the account
- `account_type` (String) Specifies the type or category of the account

### Optional

- `account_alias` (String) Required for AWS IAM accounts
- `account_expiration_date` (String) The expiration date of the account (format: DD/MM/YYYY)
- `distinguished_name` (String) Required for LDAP domain accounts
- `domain_name` (String) Required for Google Workspace accounts
- `folder_id` (Number) The ID of the folder where the account is stored
- `ipaddress` (String) The IP address of the account (if applicable)
- `notes` (String) Additional notes related to the account
- `password` (String) The password associated with the account
- `personal_account` (Boolean) Indicates whether the account is personal (true/false)
- `tags` (String) Tags associated with the account

### Read-Only

- `id` (Number) Unique identifier of the created account in Securden
- `message` (String) Response message indicating the result of the operation
