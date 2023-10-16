// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure localosProvider satisfies various provider interfaces.
var _ provider.Provider = &LocalOsProvider{}

// LocalOsProvider defines the provider implementation.
type LocalOsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// LocalOsProviderModel describes the provider data model.
type LocalOsProviderModel struct {
}

func (p *LocalOsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "localos"
	resp.Version = p.version
}

func (p *LocalOsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
The local-os provider provides information in the form of data sources about
the operating system and environment of the machine
on which you are running terraform.

In certain situations it can be useful to know if your configuration is running on Windows or not,
especically for storing locally created artifacts such as key pairs in the appropriate directories.
`,
	}
}

func (p *LocalOsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data LocalOsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *LocalOsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *LocalOsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOsInfoDataSource,
		NewFoldersDataSource,
		NewPublicIPDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LocalOsProvider{
			version: version,
		}
	}
}
