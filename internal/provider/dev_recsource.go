package provider

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/client"

	devops_resource "github.com/liatrio/devops-bootcamp/examples/ch7/devops-resources"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &devResource{}
	_ resource.ResourceWithConfigure   = &devResource{}
	_ resource.ResourceWithImportState = &devResource{}
)

// NewDevResource is a helper function to simplify the provider implementation.
func NewDevResource() resource.Resource {
	return &devResource{}
}

// devResource is the resource implementation.
type devResource struct {
	client *client.Client
}

// devResourceModel maps dev schema data.
// devModel maps dev schema data.
type devResourceModel struct {
	Name        types.String     `tfsdk:"name"`
	Id          types.String     `tfsdk:"id"`
	Engineers   []*engineerModel `tfsdk:"engineers"`
	LastUpdated types.String     `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *devResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dev_resource"
}

// Schema defines the schema for the resource.
func (r *devResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"engineers": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
							// PlanModifiers: []planmodifier.String{
							// 	stringplanmodifier.UseStateForUnknown(),
							// },
						},
						"id": schema.StringAttribute{
							Required: true,
						},
						"email": schema.StringAttribute{
							Computed: true,
							// PlanModifiers: []planmodifier.String{
							// 	stringplanmodifier.UseStateForUnknown(),
							// },
						},
					},
				},
			},
		},
	}
}

// Create a new resource.
func (r *devResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	log.Printf("Debug: Create request: %v", req)

	// Retrieve values from plan
	var plan devResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		log.Printf("Error: %v", resp.Diagnostics)
		return
	}

	var devObject devops_resource.Dev
	devObject.Name = plan.Name.ValueString()
	devObject.Id = plan.Id.ValueString()

	// Print to standard logger; this will appear if TF_LOG=DEBUG is set
	log.Printf("Debug: Dev Object: %#v", devObject)

	// Create new dev
	dev, err := r.client.CreateDev(devObject)
	if err != nil {
		log.Printf("Error: %v", err)
		resp.Diagnostics.AddError(
			"Error creating dev",
			"Could not create dev, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Name = types.StringValue(dev.Name)
	plan.Id = types.StringValue(dev.Id)
	for index, engineer := range plan.Engineers {
		ID := strings.Trim(engineer.Id.String(), "\"")

		eng, err := r.client.GetEngineer(ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error sending get request to devops-bootcamp api",
				"Could not read engineer Id "+engineer.Id.String()+": "+err.Error(),
			)
			return
		}
		err = r.client.AddEngToDev(dev.Id, eng.Id)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error sending get request to devops-bootcamp api",
				"Could not add engineer Id "+engineer.Id.String()+" to Dev "+dev.Id+": "+err.Error(),
			)
			return
		}
		plan.Engineers[index].Name = types.StringValue(eng.Name)
		plan.Engineers[index].Id = types.StringValue(eng.Id)
		plan.Engineers[index].Email = types.StringValue(eng.Email)
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		log.Printf("Error: %v", resp.Diagnostics)
		return
	}
}

// Read resource information.
func (r *devResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	log.Printf("Debug: Read request: %v", req)
	// Get current state
	var state devResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		log.Printf("Error: %v", resp.Diagnostics)
		return
	}

	// Get refreshed dev value from HashiCups
	dev, err := r.client.GetDev(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error sending get request to devops-bootcamp api",
			"Could not read dev Id "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state.Name = types.StringValue(dev.Name)
	state.Id = types.StringValue(dev.Id)
	for index, engineer := range dev.Engineers {
		state.Engineers[index].Name = types.StringValue(engineer.Name)
		state.Engineers[index].Id = types.StringValue(engineer.Id)
		state.Engineers[index].Email = types.StringValue(engineer.Email)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *devResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	log.Printf("Debug: Update request: %v", req)
	// Retrieve values from plan
	var plan devResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		log.Printf("Error: %v", resp.Diagnostics)
		return
	}

	var dev devops_resource.Dev
	// Update dev
	dev.Name = plan.Name.ValueString()
	dev.Id = plan.Id.ValueString()
	for _, engineer := range plan.Engineers {
		ID := strings.Trim(engineer.Id.String(), "\"")
		eng, err := r.client.GetEngineer(ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error sending get request to devops-bootcamp api",
				"Could not read engineer Id "+engineer.Id.String()+": "+err.Error(),
			)
			return
		}
		dev.Engineers = append(dev.Engineers, eng)

	}

	devObj, err := r.client.UpdateDev(dev)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating dev",
			"Could not update dev, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Name = types.StringValue(devObj.Name)
	plan.Id = types.StringValue(devObj.Id)
	for index, eng := range devObj.Engineers {
		plan.Engineers[index].Name = types.StringValue(eng.Name)
		plan.Engineers[index].Id = types.StringValue(eng.Id)
		plan.Engineers[index].Email = types.StringValue(eng.Email)
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		log.Printf("Error: %v", resp.Diagnostics)
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *devResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state devResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing dev
	err := r.client.DeleteDev(state.Id.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting dev",
			"Could not delete dev Id "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *devResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *devResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
