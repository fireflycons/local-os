// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
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
	client *http.Client
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

	secondaries := make([]NICModel, 0, 4)
	data.Primary = basetypes.NewObjectNull(nicAttributeTypes())
	nics, err := privateip.GetLocalIP4Interfaces(true)

	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "This is an error with the provider.")
	}

	for _, nic := range nics {
		if nic.IsPrimary {
			// Populate "primary" object attribute
			resp.Diagnostics.Append(tfsdk.ValueFrom(ctx, nicToNICModel(nic), types.ObjectType{
				AttrTypes: nicAttributeTypes(),
			}, &data.Primary)...)
		} else {
			secondaries = append(secondaries, nicToNICModel(nic))
		}
	}

	// Populate "secondaries" list of objects
	resp.Diagnostics.Append(tfsdk.ValueFrom(ctx, secondaries, types.ListType{
		ElemType: types.ObjectType{
			AttrTypes: nicAttributeTypes(),
		},
	}, &data.Secondaries)...)

	if data.Primary.IsNull() {
		resp.Diagnostics.AddError("Missing primary NIC", "No local NIC could be found that routes to a default gateway")
	} else {
		// Resource ID is combination of primary NIC IP and its network
		ipValue, ok := data.Primary.Attributes()["ip"].(basetypes.StringValue)
		if !ok {
			resp.Diagnostics.AddError("Unable to retrieve primary NIC IP as StringValue.", "This is an error with the provider.")
			return
		}
		networkValue, ok := data.Primary.Attributes()["network"].(basetypes.StringValue)
		if !ok {
			resp.Diagnostics.AddError("Unable to retrieve primary NIC network as StringValue.", "This is an error with the provider.")
			return
		}
		data.Id = types.StringValue(ipValue.ValueString() + "_" + strings.ReplaceAll(networkValue.ValueString(), "/", "_"))
	}

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
