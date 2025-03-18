package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &AddAccount{}

func add_account() datasource.DataSource {
	return &AddAccount{}
}

type AddAccount struct {
	client *http.Client
}

type AddAccountModel struct {
	AccountName           types.String `tfsdk:"account_name"`
	AccountTitle          types.String `tfsdk:"account_title"`
	AccountType           types.String `tfsdk:"account_type"`
	IPAddress             types.String `tfsdk:"ipaddress"`
	Notes                 types.String `tfsdk:"notes"`
	Tags                  types.String `tfsdk:"tags"`
	PersonalAccount       types.Bool   `tfsdk:"personal_account"`
	FolderID              types.Int64  `tfsdk:"folder_id"`
	Password              types.String `tfsdk:"password"`
	AccountExpirationDate types.String `tfsdk:"account_expiration_date"`
	DistinguishedName     types.String `tfsdk:"distinguished_name"`
	AccountAlias          types.String `tfsdk:"account_alias"`
	DomainName            types.String `tfsdk:"domain_name"`
	Message               types.String `tfsdk:"message"`
	ID                    types.Int64  `tfsdk:"id"`
}

func (d *AddAccount) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_add_account"
}

func (d *AddAccount) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Securden data source",

		Attributes: map[string]schema.Attribute{
			"account_title": schema.StringAttribute{
				MarkdownDescription: "Title of the account",
				Required:            true,
			},
			"account_name": schema.StringAttribute{
				MarkdownDescription: "Name of the account",
				Required:            true,
			},
			"account_type": schema.StringAttribute{
				MarkdownDescription: "Type of the account",
				Required:            true,
			},
			"ipaddress": schema.StringAttribute{
				MarkdownDescription: "IP Address of the account",
				Optional:            true,
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Notes of the account",
				Optional:            true,
			},
			"tags": schema.StringAttribute{
				MarkdownDescription: "Tags of the account",
				Optional:            true,
			},
			"personal_account": schema.BoolAttribute{
				MarkdownDescription: "Personal account flag (Boolean)",
				Optional:            true,
			},
			"folder_id": schema.Int64Attribute{
				MarkdownDescription: "Folder ID in which the account to be added (Long Integer)",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password of the account",
				Optional:            true,
			},
			"account_expiration_date": schema.StringAttribute{
				MarkdownDescription: "Expiration date of the account (Date Format - DD/MM/YYYY)",
				Optional:            true,
			},
			"distinguished_name": schema.StringAttribute{
				MarkdownDescription: "Required for LDAP domain accounts",
				Optional:            true,
			},
			"account_alias": schema.StringAttribute{
				MarkdownDescription: "Required for AWS IAM accounts",
				Optional:            true,
			},
			"domain_name": schema.StringAttribute{
				MarkdownDescription: "Required for Google Workspace accounts",
				Optional:            true,
			},
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"message": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *AddAccount) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AddAccount) Create(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var account AddAccountModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &account)...)
	params := make(map[string]any)
	setParam(params, "account_name", account.AccountName)
	setParam(params, "account_title", account.AccountTitle)
	setParam(params, "account_type", account.AccountType)
	setParam(params, "ipaddress", account.IPAddress)
	setParam(params, "notes", account.Notes)
	setParam(params, "tags", account.Tags)
	setParam(params, "personal_account", account.PersonalAccount)
	setParam(params, "folder_id", account.FolderID)
	setParam(params, "password", account.Password)
	setParam(params, "account_expiration_date", account.AccountExpirationDate)
	setParam(params, "distinguished_name", account.DistinguishedName)
	setParam(params, "account_alias", account.AccountAlias)
	setParam(params, "domain_name", account.DomainName)
	added_account, code, message := add_account_function(ctx, params)
	if code != 200 && code != 0 {
		resp.Diagnostics.AddWarning(fmt.Sprintf("%d - %s", code, message), "")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &added_account)...)
}

func (d *AddAccount) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var account AddAccountModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &account)...)
	params := make(map[string]any)
	setParam(params, "account_name", account.AccountName)
	setParam(params, "account_title", account.AccountTitle)
	setParam(params, "account_type", account.AccountType)
	setParam(params, "ipaddress", account.IPAddress)
	setParam(params, "notes", account.Notes)
	setParam(params, "tags", account.Tags)
	setParam(params, "personal_account", account.PersonalAccount)
	setParam(params, "folder_id", account.FolderID)
	setParam(params, "password", account.Password)
	setParam(params, "account_expiration_date", account.AccountExpirationDate)
	setParam(params, "distinguished_name", account.DistinguishedName)
	setParam(params, "account_alias", account.AccountAlias)
	setParam(params, "domain_name", account.DomainName)
	added_account, code, message := add_account_function(ctx, params)
	if code != 200 && code != 0 {
		resp.Diagnostics.AddWarning(fmt.Sprintf("%d - %s", code, message), "")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &added_account)...)
}
