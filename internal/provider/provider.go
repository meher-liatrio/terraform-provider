// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-provider-scaffolding-framework/client"

	// devops_resource "github.com/liatrio/devops-bootcamp/examples/ch7/devops-resources"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Go structs for converting Go to json
// var engineer devops_resource.Engineer

// Ensure devopsBootcampProvider satisfies various provider interfaces.
var _ provider.Provider = &devopsBootcampProvider{}
var _ provider.ProviderWithFunctions = &devopsBootcampProvider{}

// devopsBootcampProvider defines the provider implementation.
type devopsBootcampProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func (p *devopsBootcampProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "devops-bootcamp"
	resp.Version = p.version
}

// devopsBootcampProviderModel maps provider schema data to a Go type.
// uses struct types with tfsdk struct field tags to map schema definitions to Go types with the actual data
type devopsBootcampProviderModel struct {
	Host types.String `tfsdk:"host"`
}

// user defines the endpoint value when declaring this provider in the TF configuration
func (p *devopsBootcampProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "Bootcamp endpoint -- host of the app!!!",
				Required:            true,
			},
		},
	}
}

func (p *devopsBootcampProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring devops-bootcamp client")

	// Retrieve provider data from configuration
	var config devopsBootcampProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown DevOps Bootcamp Host",
			"The provider cannot create the DevOps Bootcamp client as there is an unknown configuration value for the DevOps Bootcamp host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HOST environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("HOST")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}
	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing DevOps Bootcamp Host",
			"The provider cannot create the DevOps Bootcamp client as there is a missing or empty value for the DevOps Bootcamp host. "+
				"Set the host value in the configuration or use the HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "devops-bootcamp_host", host)

	tflog.Debug(ctx, "Creating devops-bootcamp client")

	// Create a new DevOps API client using the configuration values
	client := client.NewClient(host)

	// Make the DevOps client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured devops-bootcamp client", map[string]any{"success": true})
}

func (p *devopsBootcampProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewEngineerResource,
	}
}

func (p *devopsBootcampProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEngineerDataSource,
	}
}

func (p *devopsBootcampProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &devopsBootcampProvider{
			version: version,
		}
	}
}
