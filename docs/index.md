# Securden - Terraform Provider

This guide will take you through the process of setting up and using Securden PAM as a provider block in Terraform to securely retrieve account credentials, keys, and secrets using APIs.

## Summary of Steps

1. Configuration in Securden  
2. Defining the Securden Provider Block  
3. Checking with Data Block  
4. Configuring the Output Block  
5. Accessing Data  
6. Adding Additional Fields  
7. Available Data Fields in the Plugin  

---

## 1. Configuration to be Done in Terraform

### a. Initialize Securden Provider
Add the following block to your `main.tf` file to initialize the Securden provider:

Refer [SecurdenDevOps/securden](https://registry.terraform.io/providers/SecurdenDevOps/securden/latest) | Terraform Registry for versions and updates.

### b. Declare Variables to Store Securden Attributes
Define two variables to store the Securden authentication token and server URL.

Initialize these variables with the correct values for your environment. For example, using environment variables:

- **Windows Command Prompt:**  
  ```sh
  set TF_VAR_authtoken=tf45....<authtoken>
  ```

- **Mac/Linux:**  
  ```sh
  export TF_VAR_authtoken=<authtoken>
  ```

---

## 2. Provider Block
Define the Securden provider block, referencing the previously declared variables:

```hcl
provider "securden" {
  authtoken = var.authtoken
  server_url = var.server_url
}
```

---

## 3. Using the Data Block
To fetch account data from Securden, use a data block. Here’s an example to fetch SSH key credentials:

```hcl
data "securden_keyvalue" "ssh" {
  account_id = "2000000002800"
}
```

> **Note:** You can retrieve specific accounts using `account_id`, `account_name`, or `account_title`. If you use `account_name`, you can also include `account_title` to ensure you don’t retrieve multiple accounts in case there are several accounts with the same name.

---

## 4. Output Block
Using the Output block to display data fetched from Securden by the plugin:

```hcl
output "ssh_password" {
  value     = data.securden_keyvalue.ssh.password
  sensitive = true
}
```

Setting the variable `sensitive = true` hides the credential value during execution, preventing the credentials from being accidentally displayed in return values or error messages.

---

## 5. Accessing Data
Here are some examples of how to access various credentials from the Securden data block:

- **For Password:** `data.securden_keyvalue.ssh.password`  
- **For PuTTY Private Key:** `data.securden_keyvalue.ssh.putty_private_key`  
- **For PuTTY Passphrase:** `data.securden_keyvalue.ssh.ppk_passphrase`  

---

## 6. Additional Fields
If your account type has additional fields in Securden, you can retrieve the value of additional fields by specifying a `key_field` and `key_value`:

```hcl
data "securden_keyvalue" "ssh" {
  account_id = "2000000002800"
  key_field  = "custom_field"
  key_value  = "field_value"
}
```

---

## 7. Available Data Fields in Plugin
Here is a list of the account attributes that can be retrieved for use in Terraform using the Securden plugin:

- `account_id`
- `account_name`
- `account_title`
- `password`
- `key_value`
- `private_key`
- `putty_private_key`
- `passphrase`
- `ppk_passphrase`
- `address`
- `client_id`
- `client_secret`
- `account_alias`
- `account_file`
- `default_database`
- `sql_server_port`
- `mysql_port`
- `oracle_sid`
- `oracle_service_name`
- `oracle_port`

> **Note:** Data can only be retrieved for the attributes that are available in the account. For any other fields, or if there is a non-existent value, a null value will be returned when the code is executed.

---

## Bulk Password Retrieval
You have the option to fetch account passwords in bulk from Securden at once using a data block.

The major difference between `securden_account` and `securden_passwords` commands is that:
- `securden_account` will raise a request each time it is called by the data block for a single account.
- `securden_passwords` will retrieve multiple account passwords in a single fetch, reducing overall time consumption.

Here’s an example of a data block used to fetch multiple account passwords:

```hcl
data "securden_passwords" "passwords" {}
```

Accounts whose passwords need to be retrieved can be called by their respective account IDs:

```hcl
data.securden_passwords.passwords["2000000002800"]
```

