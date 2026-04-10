---
layout: "mysql"
page_title: "Provider: TiDB"
sidebar_current: "docs-mysql-index"
description: |-
  A Terraform provider for managing TiDB users, roles, grants, and TiDB-specific configuration.
---

# TiDB Provider

The provider in this repository is a TiDB-oriented Terraform provider built on top of the MySQL provider model.

It is intended for teams that manage TiDB access declaratively with Terraform and need stable lifecycle behavior for users and roles.

## What this provider manages

The provider is primarily used for:

- TiDB users
- role grants
- default roles
- TiDB configuration resources

Resource names still use the historical `mysql_*` prefix for backward compatibility, but the provider is documented and maintained with TiDB as a first-class target.

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

## Common Role Grant Pattern

```hcl
resource "mysql_user" "readonly_user" {
  user       = "alice"
  host       = "%"
  tls_option = "X509"
}

resource "mysql_grant" "readonly_role" {
  user  = mysql_user.readonly_user.user
  host  = mysql_user.readonly_user.host
  roles = ["readonly_role"]
}
```

## Why this fork exists

The upstream MySQL-compatible provider is a strong base, but TiDB role lifecycle handling requires more careful quoting and normalization than generic MySQL-oriented assumptions provide. This provider improves predictability for:

- `GRANT role TO user`
- `REVOKE role FROM user`
- `ALTER USER ... DEFAULT ROLE`
- mixed role identifier formats in configuration and state, for example:
  - `readonly_role`
  - `readonly_role@%`
  - `'readonly_role'@'%'`

## Notes

- Resource names currently retain the `mysql_*` prefix for backward compatibility.
- Legacy quoted role identifiers such as `"'readonly_role'@'%'"` are normalized where supported.
- For TiDB user management, the most relevant resources are `mysql_user`, `mysql_grant`, `mysql_default_roles`, `mysql_role`, and `mysql_ti_config`.
