// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"
	"os"

	"github.com/Serviceware/terraform-provider-swp/internal/aipe"
	"github.com/Serviceware/terraform-provider-swp/internal/authenticator"
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

type AIPEProvider struct {
	version string
}

type AIPEProviderModel struct {
	ApplicationUsername   types.String `tfsdk:"application_username"`
	ApplicationPassword   types.String `tfsdk:"application_password"`
	AuthenticatorRealmURL types.String `tfsdk:"authenticator_realm_url"`
	AIPEURL               types.String `tfsdk:"aipe_url"`
}

func (p *AIPEProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "swp"
	resp.Version = p.version
}

func (p *AIPEProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"application_username": schema.StringAttribute{
				MarkdownDescription: "Username for AIPE from user management",
				Optional:            true,
			},
			"application_password": schema.StringAttribute{
				MarkdownDescription: "Password for AIPE from user management",
				Optional:            true,
				Sensitive:           true,
			},

			"authenticator_realm_url": schema.StringAttribute{
				MarkdownDescription: "URL of the Authenticatopr realm",
				Optional:            true,
			},
			"aipe_url": schema.StringAttribute{
				MarkdownDescription: "URL of the AIPE",
				Optional:            true,
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

	client := http.DefaultClient

	applicationUsername := os.Getenv("SWP_APPLICATION_USER_USERNAME")
	applicationPassword := os.Getenv("SWP_APPLICATION_USER_PASSWORD")
	authenticatorRealmURL := os.Getenv("SWP_AUTHENTICATOR_URL")
	aipeURL := os.Getenv("SWP_AIPE_URL")

	if data.ApplicationUsername.ValueString() != "" {
		applicationUsername = data.ApplicationUsername.ValueString()
	}

	if data.ApplicationPassword.ValueString() != "" {
		applicationPassword = data.ApplicationPassword.ValueString()
	}

	if data.AuthenticatorRealmURL.ValueString() != "" {
		authenticatorRealmURL = data.AuthenticatorRealmURL.ValueString()
	}

	if data.AIPEURL.ValueString() != "" {
		aipeURL = data.AIPEURL.ValueString()
	}

	if applicationUsername == "" {
		resp.Diagnostics.AddError("application_username", "application_username is required")
	}

	if applicationPassword == "" {
		resp.Diagnostics.AddError("application_password", "application_password is required")
	}

	if authenticatorRealmURL == "" {
		resp.Diagnostics.AddError("authenticator_realm_url", "authenticator_realm_url is required")
	}

	if aipeURL == "" {
		resp.Diagnostics.AddError("aipe_url", "aipe_url is required")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	authenticatorClient := authenticator.AuthenticatorClient{
		Client:              client,
		ApplicationUsername: applicationUsername,
		ApplicationPassword: applicationPassword,
		URL:                 authenticatorRealmURL,
	}

	token, err := authenticatorClient.Authenticate(ctx)

	if err != nil {
		resp.Diagnostics.AddError("authenticator", "failed to authenticate: "+err.Error())
		return
	}

	aipeClient := aipe.AIPEClient{
		HTTPClient: client,
		URL:        aipeURL,
		OIDCToken:  token,
	}

	tflog.Info(ctx, "Successfully configured AIPE provider")
	resp.DataSourceData = &aipeClient
	resp.ResourceData = &aipeClient
}

func (p *AIPEProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDataObjectResource,
		NewDataObjectLinkResource,
	}
}

func (p *AIPEProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDataObjectDataSource,
	}
}

func (p *AIPEProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AIPEProvider{
			version: version,
		}
	}
}
