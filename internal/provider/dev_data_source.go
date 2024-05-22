package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/client"
)

// var dev devops_resource.Dev

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &devDataSource{}
	_ datasource.DataSourceWithConfigure = &devDataSource{}
)

// NewDevDataSource is a helper function to simplify the provider implementation.
func NewDevDataSource() datasource.DataSource {
	return &devDataSource{}
}

// devDataSource is the data source implementation.
type devDataSource struct {
	client *client.Client
}

// devDataSourceModel maps the data source schema data.
type devDataSourceModel struct {
	Devs []devModel `tfsdk:"devs"`
}

// hard coding for readability -- can also be imported via devops_resource
// devModel maps dev schema data.
type devModel struct {
	Name        types.String     `tfsdk:"name"`
	Id          types.String     `tfsdk:"id"`
	Engineers   []*engineerModel `tfsdk:"engineers"`
	LastUpdated types.String     `tfsdk:"last_updated"`
}

// Metadata returns the data source type name.
func (d *devDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devs"
}

// Schema defines the schema for the data source.
func (d *devDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Devs data source",

		Attributes: map[string]schema.Attribute{
			"devs": schema.ListNestedAttribute{
				MarkdownDescription: "Dev attribute",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Dev id computed",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Dev name required",
							Required:            true,
						},
						"engineers": schema.ListNestedAttribute{
							MarkdownDescription: "List of Engineers optional",
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										MarkdownDescription: "Engineer id required",
										Required:            true,
									},
									"name": schema.StringAttribute{
										MarkdownDescription: "Engineer name computed",
										Computed:            true,
										// PlanModifiers: []planmodifier.String{
										// 	stringplanmodifier.UseStateForUnknown(),
										// },
									},
									"email": schema.StringAttribute{
										MarkdownDescription: "Engineer email computed",
										Computed:            true,
										// PlanModifiers: []planmodifier.String{
										// 	stringplanmodifier.UseStateForUnknown(),
										// },
									},
								},
							},
						},
						"last_updated": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *devDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state devDataSourceModel

	devs, err := d.client.GetDevs()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read DevOps Dev",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, dev := range devs {
		tempDev := devModel{
			Id:   types.StringValue(dev.Id),
			Name: types.StringValue(dev.Name),
		}
		for _, engineer := range dev.Engineers {
			tempEngineer := engineerModel{
				Id:    types.StringValue(engineer.Id),
				Name:  types.StringValue(engineer.Name),
				Email: types.StringValue(engineer.Email),
			}
			tempDev.Engineers = append(tempDev.Engineers, &tempEngineer)
		}
		state.Devs = append(state.Devs, tempDev)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *devDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
