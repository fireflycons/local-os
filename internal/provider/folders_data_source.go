// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/fireflycons/local-os/internal/provider/specialfolder"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &FoldersDataSource{}

func NewFoldersDataSource() datasource.DataSource {
	return &FoldersDataSource{}
}

// FoldersDataSource defines the data source implementation.
type FoldersDataSource struct {
}

// FoldersDataSourceModel describes the data source data model.
type FoldersDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Home types.String `tfsdk:"home"`
	SSH  types.String `tfsdk:"ssh"`
}

func (d *FoldersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folders"
}

func (d *FoldersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The `folders` data source gets information about special folders on the machine that is running terraform.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource identifier",
				Computed:            true,
			},
			"home": schema.StringAttribute{
				MarkdownDescription: "Absolute path to user's home directory",
				Computed:            true,
			},
			"ssh": schema.StringAttribute{
				MarkdownDescription: "Absolute path to user's SSH keys directory",
				Computed:            true,
			},
		},
	}
}

func (d *FoldersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Nothing to configure
}

func (d *FoldersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FoldersDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	f := specialfolder.NewSpecialFolder()
	data.Id = types.StringValue(f.ID())
	data.Home = types.StringValue(f.Home())
	data.SSH = types.StringValue(f.SSH())

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "Read folders data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
