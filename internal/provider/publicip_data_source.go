// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	url1 "net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Use this resource as it simply returns the IP in plain text,
// without being a full HTML document.
const amazonCheckIp = "https://checkip.amazonaws.com"

var checkIpUrl, _ = url1.Parse(amazonCheckIp)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PublicIPDataSource{}

func NewPublicIPDataSource() datasource.DataSource {
	return &PublicIPDataSource{}
}

// PublicIPDataSource defines the data source implementation.
type PublicIPDataSource struct {
	client *http.Client
}

// PublicIPDataSourceModel describes the data source data model.
type PublicIPDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Cidr types.String `tfsdk:"cidr"`
	IP   types.String `tfsdk:"ip"`
}

func (d *PublicIPDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_ip"
}

func (d *PublicIPDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "`public_ip` data source gets the public IP of the machine that reads the data source as a CIDR.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource identifier",
				Computed:            true,
			},
			"cidr": schema.StringAttribute{
				MarkdownDescription: `/32 public IP CIDR of machine running terraform,


    This is useful when deploying e.g. test infrastructure for which you want to only grant access to your own workstation.
    You can use this to set up firewalls, cloud security groups etc.

    Will be empty string if the internet is not accessible from the caller.
`,
				Computed:  true,
				Sensitive: true,
			},
			"ip": schema.StringAttribute{
				MarkdownDescription: `Public IP of machine running terraform.`,
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (d *PublicIPDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PublicIPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PublicIPDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpReq := &http.Request{
		Method: "GET",
		URL:    checkIpUrl,
	}

	httpResp, err := d.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read %s, got error: %s", amazonCheckIp, err))
		data.Cidr = types.StringValue("")
		return
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, httpResp.Body)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read response from %s, got error: %s", amazonCheckIp, err))
		data.Cidr = types.StringValue("")
		return
	}

	ip := strings.TrimSpace(buf.String())
	data.Id = types.StringValue(amazonCheckIp)
	data.Cidr = types.StringValue(ip + "/32")
	data.IP = types.StringValue(ip)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "Read public_ip data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
