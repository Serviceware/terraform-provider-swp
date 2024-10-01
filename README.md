# Terraform Provider Serviceware Platform

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Currently, I have not published this provider to the registry yet in case we need to make
bigger changes to the resources.

As such, you need to install it via dev-overrides.

1. Build the provider

2. Setup dev-dependencies. Edit your `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
     "swp" = "path-to-checkout"
  }

  direct {}
}
```

3. No `terraform init` necessary, `terraform plan` and `terraform apply` pickup the dev-dependency directly.

4. Setup a provider (see the `docs/` folder) and start buildling data objects.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

