package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/client"
)

// var engineer devops_resource.Engineer

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &engineerDataSource{}
	_ datasource.DataSourceWithConfigure = &engineerDataSource{}
)

// NewEngineerDataSource is a helper function to simplify the provider implementation.
func NewEngineerDataSource() datasource.DataSource {
	return &engineerDataSource{}
}

// engineerDataSource is the data source implementation.
type engineerDataSource struct {
	client *client.Client
}

// engineerDataSourceModel maps the data source schema data.
type engineerDataSourceModel struct {
	Engineer []engineerModel `tfsdk:"engineer"`
	// implemented to allow testing framework to interact with this data source
	// placeholder value is applied in Read method
	// ID types.String `tfsdk:"id"`
}

// hard coding for readability -- can also be imported via devops_resource
// engineerModel maps engineer schema data.
type engineerModel struct {
	Name  types.String `tfsdk:"name"`
	Id    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
}

// Metadata returns the data source type name.
func (d *engineerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer"
}

// Schema defines the schema for the data source.
func (d *engineerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"engineer": schema.ListNestedAttribute{
				MarkdownDescription: "Engineer attribute",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Engineer Name required",
							Required:            true,
						},
						"id": schema.StringAttribute{
							MarkdownDescription: "Engineer ID computed",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "Engineer Email required",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *engineerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state engineerDataSourceModel

	engineers, err := d.client.GetEngineers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HashiCups Engineer",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, engineer := range engineers {
		engineerState := engineerModel{
			Name:  types.StringValue(engineer.Name),
			Id:    types.StringValue(engineer.Id),
			Email: types.StringValue(engineer.Email),
		}

		state.Engineer = append(state.Engineer, engineerState)
	}

	// // data source needs an id attr for testing framework to perform tests
	// // this data source doesn't actually have an id so we just use a placeholder
	// state.ID = types.StringValue("placeholder")

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *engineerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
