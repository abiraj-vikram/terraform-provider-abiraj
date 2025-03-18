package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &securdenProvider{}
var _ provider.ProviderWithFunctions = &securdenProvider{}

type securdenProvider struct {
	version string
}

var SecurdenAuthToken string
var SecurdenServerURL string
var SecurdenOrg string
var SecurdenCertificate string
var PluginVersion string

type securdenProviderModel struct {
	ServerURL   types.String `tfsdk:"server_url"`
	AuthToken   types.String `tfsdk:"authtoken"`
	Certificate types.String `tfsdk:"certificate"`
}

func (p *securdenProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "securden"
	resp.Version = p.version
}

func (p *securdenProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				Required: true,
			},
			"authtoken": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"certificate": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Securden Server SSL Certificate",
			},
		},
	}
}

func (p *securdenProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config securdenProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	SecurdenServerURL = config.ServerURL.ValueString()
	SecurdenCertificate = config.Certificate.ValueString()
	SecurdenAuthToken = config.AuthToken.ValueString()
	PluginVersion = p.version
}

func (p *securdenProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *securdenProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		account,
		accounts,
		accounts_passwords,
		add_account,
		edit_account,
		delete_accounts,
	}
}

func (p *securdenProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{}
}

func Provider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &securdenProvider{
			version: version,
		}
	}
}
