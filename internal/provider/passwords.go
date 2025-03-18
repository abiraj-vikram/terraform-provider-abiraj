package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &AccountsPasswords{}

func accounts_passwords() datasource.DataSource {
	return &AccountsPasswords{}
}

type AccountsPasswords struct {
	client *http.Client
}

type AccountsPasswordsModel struct {
	AccountIDs []types.String `tfsdk:"account_ids"`
	Passwords  types.Map      `tfsdk:"passwords"`
}

func (d *AccountsPasswords) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_passwords"
}

func (d *AccountsPasswords) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Securden Accounts Passwords",

		Attributes: map[string]schema.Attribute{
			"account_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of account ids",
				Required:            true,
			},
			"passwords": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *AccountsPasswords) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*http.Client)
	if !ok {
		resp.Diagnostics.AddWarning(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *AccountsPasswords) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if PluginVersion != "1.0.0" {
		resp.Diagnostics.AddWarning("The feature is no more supported", "")
		return
	}
	var account AccountsPasswordsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &account)...)

	var accountIDs []string
	for _, id := range account.AccountIDs {
		accountIDs = append(accountIDs, id.ValueString())
	}

	passwords, code, message := get_passwords(ctx, accountIDs)
	if code != 200 {
		resp.Diagnostics.AddWarning(fmt.Sprintf("%d - %s", code, message), "")
		return
	}

	account.Passwords = passwords

	resp.Diagnostics.Append(resp.State.Set(ctx, &account)...)
}

func (d *AccountsPasswords) Create(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if PluginVersion != "1.0.0" {
		resp.Diagnostics.AddWarning("The feature is no more supported", "")
		return
	}
	var account AccountsPasswordsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &account)...)

	var accountIDs []string
	for _, id := range account.AccountIDs {
		accountIDs = append(accountIDs, id.ValueString())
	}

	passwords, code, message := get_passwords(ctx, accountIDs)
	if code != 200 {
		resp.Diagnostics.AddWarning(fmt.Sprintf("%d - %s", code, message), "")
		return
	}

	account.Passwords = passwords

	resp.Diagnostics.Append(resp.State.Set(ctx, &account)...)
}
