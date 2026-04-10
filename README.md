# Terraform Provider TiDB

Terraform provider for managing TiDB users, grants, roles, and related access objects.

This repository is based on the [`petoju/terraform-provider-mysql`](https://github.com/petoju/terraform-provider-mysql) fork and is being adapted for reliable TiDB lifecycle management.

## Why this fork exists

The upstream MySQL-compatible provider is a strong base, but TiDB-specific issues around role lifecycle management require additional care, especially for:

- `GRANT role TO user`
- `REVOKE role FROM user`
- default roles
- role identifiers stored in mixed formats such as:
  - `readonly_role`
  - `readonly_role@%`
  - `'readonly_role'@'%'`

This fork exists so we can make TiDB behavior explicit and predictable instead of relying on MySQL-oriented quoting/parsing assumptions.

## Current focus

The provider is being adapted to support TiDB user management as a first-class Terraform workflow, with emphasis on:

- certificate-based users
- declarative role membership
- stable read/plan/apply behavior
- compatibility with existing legacy state formats

## What has already been fixed

The first TiDB-specific patch in this fork fixes role identifier normalization for:

- `mysql_grant`
- `mysql_default_roles`
- `mysql_role`

The provider now canonicalizes role identifiers before generating SQL, so mixed role formats are converted into a single stable form before `GRANT`, `REVOKE`, `ALTER USER ... DEFAULT ROLE`, and role operations are executed.

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

## Planned direction

Short-term goal:

- make TiDB user and role lifecycle reliable in existing modules

Mid-term goal:

- provide a clean Terraform-native way to manage TiDB users in a desired-state workflow

Possible future direction:

- introduce explicitly TiDB-oriented resources if MySQL compatibility becomes a blocker for clean semantics

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) or OpenTofu
- [Go](https://go.dev/doc/install)

## Local development

### Build

```bash
go build ./...
```

### Run targeted tests

```bash
go test ./mysql -run 'TestNormalizeRoleIdentifier|TestRoleIdentifierSQL|TestRoleGrantSQLStatements'
```

### Format

```bash
gofmt -w .
```

## Releasing

Releases are published from Git tags matching `v*` using GitHub Actions and GoReleaser.

Before the first release, configure these GitHub Actions secrets:

- `GPG_PRIVATE_KEY`
- `GPG_PASSPHRASE`
- `GPG_FINGERPRINT`

Then create and push a tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow will build provider artifacts, sign checksums, and create a GitHub Release suitable for Terraform Registry ingestion.

## Using the provider

Until the provider is published to Terraform Registry, local development can use a local build or override installation.

After publication, the intended usage will look like:

```hcl
terraform {
  required_providers {
    tidb = {
      source  = "mf-troy/tidb"
      version = "~> 0.1"
    }
  }
}
```

## Upstream relationship

- this repository is a fork of `petoju/terraform-provider-mysql`
- upstream is kept as a Git remote for easier rebasing/cherry-picking
- TiDB-specific behavior is implemented here intentionally rather than handled ad hoc outside the provider

## Repository goals

This repository should become:

- the public home of the TiDB provider
- the source of truth for TiDB-specific provider fixes
- the provider used instead of generic MySQL behavior where TiDB semantics differ
