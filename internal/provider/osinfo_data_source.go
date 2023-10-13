// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"runtime"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &OsInfoDataSource{}

func NewOsInfoDataSource() datasource.DataSource {
	return &OsInfoDataSource{}
}

// OsInfoDataSource defines the data source implementation.
type OsInfoDataSource struct {
}

// OsInfoDataSourceModel describes the data source data model.
type OsInfoDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Arch    types.String `tfsdk:"arch"`
	Windows types.Bool   `tfsdk:"is_windows"`
}

func (d *OsInfoDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_info"
}

func (d *OsInfoDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The `info` data source gets information about the operating system of the machine that is running terraform.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource identifier",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "OS Name, e.g. \"linux\"",
				Computed:            true,
			},
			"arch": schema.StringAttribute{
				MarkdownDescription: "OS Architecture, e.g. \"amd64\"",
				Computed:            true,
			},
			"is_windows": schema.BoolAttribute{
				MarkdownDescription: "Utility attribute to quickly determine windows/not windows. Other supported OS are assumed to follow POSIX semantics.",
				Computed:            true,
			},
		},
	}
}

func (d *OsInfoDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Nothing to configure
}

func (d *OsInfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OsInfoDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	os := runtime.GOOS
	arch := runtime.GOARCH

	data.Id = types.StringValue(os + "/" + arch)
	data.Name = types.StringValue(os)
	data.Arch = types.StringValue(arch)
	data.Windows = types.BoolValue(os == "windows")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "Read osinfo data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
