// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/fireflycons/terraform-provider-localos/internal/helpers/privateip"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PrivateIPDataSource{}

func NewPrivateIPDataSource() datasource.DataSource {
	return &PrivateIPDataSource{}
}

// PrivateIPDataSource defines the data source implementation.
type PrivateIPDataSource struct {
	localInterfaces privateip.LocalInterfaces
}

// PrivateIPDataSourceModel describes the data source data model.
type PrivateIPDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Primary     types.Object `tfsdk:"primary"`
	Secondaries types.List   `tfsdk:"secondaries"` //< NICModel
}

type NICModel struct {
	Name    types.String `tfsdk:"name"`
	Ip      types.String `tfsdk:"ip"`
	Cidr    types.String `tfsdk:"cidr"`
	Network types.String `tfsdk:"network"`
}

func (d *PrivateIPDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_ip"
}

func (d *PrivateIPDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "`private_ip` data source gets information about all the network interfaces attached to the machine that reads the data source. " +
			"All discovered interfaces except loopback are returned, as loopback is well known. ",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource identifier",
				Computed:            true,
			},
			"primary": schema.ObjectAttribute{
				MarkdownDescription: "Primary NIC (routes to default gateway)",
				Computed:            true,
				AttributeTypes:      nicAttributeTypes(),
			},
			"secondaries": schema.ListAttribute{
				MarkdownDescription: "All other NICs",
				Computed:            true,
				ElementType: types.ObjectType{
					AttrTypes: nicAttributeTypes(),
				},
			},
		},
	}
}

func (d *PrivateIPDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	configData, ok := req.ProviderData.(ConfigurationData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.localInterfaces = configData.localInterfaces
}

func nicAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":    types.StringType,
		"ip":      types.StringType,
		"cidr":    types.StringType,
		"network": types.StringType,
	}
}

func (d *PrivateIPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PrivateIPDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read local interfaces
	if err := d.localInterfaces.ScanInterfaces(); err != nil {
		resp.Diagnostics.AddError(err.Error(), "This is an error with the provider.")
	}

	// Populate "primary"
	p := d.localInterfaces.GetPrimary()

	if p == nil {
		data.Primary = basetypes.NewObjectNull(nicAttributeTypes())
		resp.Diagnostics.AddWarning("No primary network interface detected", d.localInterfaces.GetPrimaryAbsentReason())
	} else {
		resp.Diagnostics.Append(tfsdk.ValueFrom(ctx, nicToNICModel(p), types.ObjectType{
			AttrTypes: nicAttributeTypes(),
		}, &data.Primary)...)
	}

	// Populate "secondaries"
	secondaries := make([]NICModel, 0, 4)

	for _, nic := range d.localInterfaces.GetSecondaries() {
		secondaries = append(secondaries, nicToNICModel(nic))
	}

	resp.Diagnostics.Append(tfsdk.ValueFrom(ctx, secondaries, types.ListType{
		ElemType: types.ObjectType{
			AttrTypes: nicAttributeTypes(),
		},
	}, &data.Secondaries)...)

	// Generate resource ID
	// Resource ID is combination of first NIC IP and its network
	ridNic := d.localInterfaces.GetFirst()

	if ridNic == nil {
		resp.Diagnostics.AddError("No local NICs found", "Either this machine has no TCP/IP interfaces, or it is an error with the provider")
		return
	}

	data.Id = types.StringValue(ridNic.Ip + "_" + strings.ReplaceAll(ridNic.Network, "/", "_"))

	if resp.Diagnostics.HasError() {
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "Read private_ip data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func nicToNICModel(nic *privateip.NIC) NICModel {
	return NICModel{
		Name:    types.StringValue(nic.Name),
		Ip:      types.StringValue(nic.Ip),
		Cidr:    types.StringValue(fmt.Sprintf("%s/32", nic.Ip)),
		Network: types.StringValue(nic.Network),
	}
}
