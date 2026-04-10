---
page_title: "mf-troy/tidb Provider"
description: |-
  Terraform provider for managing TiDB users, roles, grants, and TiDB-specific configuration with MySQL-compatible resources.
---

# mf-troy/tidb Provider

The `mf-troy/tidb` provider is a TiDB-oriented Terraform provider built on top of the MySQL provider model.

It is designed for teams that need to manage TiDB access declaratively with Terraform, especially when:

- database users are provisioned centrally
- X509 / certificate-based access is required
- roles must be granted and revoked predictably
- legacy role formats already exist in Terraform state

## What this provider is for

The provider focuses on reliable TiDB user lifecycle management:

- create and update users
- enforce TLS / X509 authentication requirements
- grant roles to users and roles
- manage default roles
- manage TiDB configuration resources already supported by the upstream codebase

Although many resources still use the historical `mysql_*` naming prefix, this provider should be treated as TiDB-first in behavior and documentation.

## Why this provider exists

The upstream MySQL-compatible provider is a strong base, but TiDB role lifecycle handling needs tighter guarantees than generic MySQL-oriented quoting and parsing provide. In particular, this fork improves stability around:

- `GRANT role TO user`
- `REVOKE role FROM user`
- `ALTER USER ... DEFAULT ROLE`
- mixed role identifier formats such as:
  - `readonly_role`
  - `readonly_role@%`
  - `'readonly_role'@'%'`

## Getting started

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

## Common role grant pattern

This is a common TiDB access pattern:

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

## Important notes

- Resource names currently retain the `mysql_*` prefix for backward compatibility.
- The provider normalizes legacy quoted role identifiers where supported.
- This provider is intended to become the stable Terraform interface for TiDB user and role management, instead of relying on generic MySQL behavior where TiDB semantics differ.

## Most relevant resources

For TiDB user management, start with:

- `mysql_user`
- `mysql_grant`
- `mysql_default_roles`
- `mysql_role`
- `mysql_ti_config`
