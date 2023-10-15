// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Arch        types.String `tfsdk:"arch"`
	Windows     types.Bool   `tfsdk:"is_windows"`
	Environment types.Map    `tfsdk:"environment"`
}

func (d *OsInfoDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_info"
}

func (d *OsInfoDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The `info` data source gets information about the operating system and environment of the machine that is running terraform.",

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
			"environment": schema.MapAttribute{
				MarkdownDescription: "Map of all environment variables",
				ElementType:         types.StringType,
				Computed:            true,
				Sensitive:           true,
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

	goos := runtime.GOOS
	arch := runtime.GOARCH
	data.Id = types.StringValue(goos + "/" + arch)
	data.Name = types.StringValue(goos)
	data.Arch = types.StringValue(arch)
	data.Windows = types.BoolValue(goos == "windows")

	env := make(map[string]attr.Value)
	for _, envvar := range os.Environ() {
		s := strings.Split(envvar, "=")
		env[s[0]] = types.StringValue(s[1])
	}

	m, diags := types.MapValue(types.StringType, env)
	data.Environment = m
	resp.Diagnostics.Append(diags...)
	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "Read osinfo data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
