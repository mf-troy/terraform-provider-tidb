---
page_title: "mf-troy/tidb Provider"
description: |-
  Terraform provider for managing TiDB users, grants, roles, and related access objects.
---

# mf-troy/tidb Provider

The `mf-troy/tidb` provider manages TiDB users, grants, roles, and related access objects.

This provider is based on the `petoju/mysql` fork and is being adapted for reliable TiDB lifecycle management.

## Why this provider exists

We use Terraform to manage a large number of TiDB database users.

While TiDB is MySQL-compatible in many areas, role lifecycle behavior requires more careful handling than the generic MySQL provider currently offers. In particular, this provider focuses on stable Terraform behavior for:

- role grants
- role revokes
- default roles
- mixed role identifier formats

## Example Usage

```hcl
terraform {
  required_providers {
    tidb = {
      source  = "mf-troy/tidb"
      version = "~> 0.1"
    }
  }
}

provider "tidb" {
  endpoint = "tidb.example.com:4000"
  username = "terraform-admin"
  password = var.tidb_admin_password
  tls      = "skip-verify"
}
```

## Example: X509 user for Teleport

```hcl
resource "mysql_user" "teleport_reader" {
  user       = "alice"
  host       = "%"
  tls_option = "X509"
}

resource "mysql_grant" "teleport_reader_role" {
  user  = mysql_user.teleport_reader.user
  host  = mysql_user.teleport_reader.host
  roles = ["teleport_reader"]
}
```

## Notes

- Resource names currently retain the `mysql_*` prefix for backward compatibility while TiDB-specific behavior is improved in the provider implementation.
- Legacy quoted role identifiers such as `"'teleport_reader'@'%'"` are normalized automatically by this provider where supported.
- The long-term goal of this provider is reliable desired-state TiDB user management from Terraform modules.
