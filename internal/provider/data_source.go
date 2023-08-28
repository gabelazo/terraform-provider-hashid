// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/speps/go-hashids"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &HashIdDataSource{}

func NewHashIdDataSource() datasource.DataSource {
	return &HashIdDataSource{}
}

// hashidDataSource defines the data source implementation.
type HashIdDataSource struct {
	client *http.Client
}

// HashIdDataSourceModel describes the data source data model.
type HashIdDataSourceModel struct {
	Alphabet    types.String `tfsdk:"alphabet"`
	MinLength   types.Int64  `tfsdk:"min_length"`
	Salt        types.String `tfsdk:"salt"`
	EncodeValue types.String `tfsdk:"encode_value"`
	HashId      types.String `tfsdk:"hash_id"`
}

func (d *HashIdDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_encode"
}

func (d *HashIdDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]schema.Attribute{
			"alphabet": schema.StringAttribute{
				MarkdownDescription: "",
				Required:            true,
			},
			"min_length": schema.Int64Attribute{
				MarkdownDescription: "",
				Required:            true,
			},
			"salt": schema.StringAttribute{
				MarkdownDescription: "salt",
				Required:            true,
			},
			"encode_value": schema.StringAttribute{
				MarkdownDescription: "string",
				Required:            true,
			},
			"hash_id": schema.StringAttribute{
				MarkdownDescription: "hashid",
				Computed:            true,
			},
		},
	}
}

func (d *HashIdDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *HashIdDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data HashIdDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hd := hashids.NewData()
	hd.Salt = data.Salt.ValueString()
	hd.Alphabet = data.Alphabet.ValueString()
	// hd.MinLength = int(data.MinLength.ValueInt64())
	h, _ := hashids.NewWithData(hd)
	hashid, _ := h.EncodeHex(hex.EncodeToString([]byte(data.EncodeValue.ValueString())))

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.HashId = types.StringValue(hashid)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
