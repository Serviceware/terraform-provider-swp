// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/Serviceware/terraform-provider-aipe/internal/aipe"
	"github.com/Serviceware/terraform-provider-aipe/internal/authenticator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &AIPEProvider{}
var _ provider.ProviderWithFunctions = &AIPEProvider{}

// AIPEProvider defines the provider implementation.
type AIPEProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type AIPEProviderModel struct {
	ApplicationUsername   types.String `tfsdk:"application_username"`
	ApplicationPassword   types.String `tfsdk:"application_password"`
	AuthenticatorRealmURL types.String `tfsdk:"authenticator_realm_url"`
	AIPEURL               types.String `tfsdk:"aipe_url"`
}

func (p *AIPEProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "aipe"
	resp.Version = p.version
}

func (p *AIPEProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"application_username": schema.StringAttribute{
				MarkdownDescription: "Username for AIPE from user management",
				Required:            true,
			},
			"application_password": schema.StringAttribute{
				MarkdownDescription: "Password for AIPE from user management",
				Required:            true,
				Sensitive:           true,
			},

			"authenticator_realm_url": schema.StringAttribute{
				MarkdownDescription: "URL of the Authenticatopr realm",
				Required:            true,
			},
			"aipe_url": schema.StringAttribute{
				MarkdownDescription: "URL of the AIPE",
				Required:            true,
			},
		},
	}
}

func (p *AIPEProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AIPEProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient

	authenticatorClient := authenticator.AuthenticatorClient{
		Client:              client,
		ApplicationUsername: data.ApplicationUsername.ValueString(),
		ApplicationPassword: data.ApplicationPassword.ValueString(),
		URL:                 data.AuthenticatorRealmURL.ValueString(),
	}

	token, err := authenticatorClient.Authenticate(ctx)

	if err != nil {
		resp.Diagnostics.AddError("authenticator", "failed to authenticate: "+err.Error())
		return
	}

	aipeClient := aipe.AIPEClient{
		HTTPClient: client,
		URL:        data.AIPEURL.ValueString(),
		OIDCToken:  token,
	}

	tflog.Info(ctx, "Successfully configured AIPE provider")
	resp.DataSourceData = &aipeClient
	resp.ResourceData = &aipeClient
}

func (p *AIPEProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *AIPEProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDataObjectDataSource,
	}
}

func (p *AIPEProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AIPEProvider{
			version: version,
		}
	}
}
