package provider

import (
	"context"
	"fmt"
	"log"
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
	_ resource.Resource                = &engineerResource{}
	_ resource.ResourceWithConfigure   = &engineerResource{}
	_ resource.ResourceWithImportState = &engineerResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewEngineerResource() resource.Resource {
	return &engineerResource{}
}

// engineerResource is the resource implementation.
type engineerResource struct {
	client *client.Client
}

// engineerResourceModel maps engineer schema data.
type engineerResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Id          types.String `tfsdk:"id"`
	Email       types.String `tfsdk:"email"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *engineerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer_resource"
}

// Schema defines the schema for the resource.
func (r *engineerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true, // Name must be provided by the user
			},
			"email": schema.StringAttribute{
				Required: true, // Email must be provided by the user
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create a new engineer resource.
func (r *engineerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	log.Printf("Debug: Create request: %v", req)
	// Retrieve values from plan
	var plan engineerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		log.Printf("Error: %v", resp.Diagnostics)
		return
	}

	var engineerObject devops_resource.Engineer
	engineerObject.Name = plan.Name.ValueString()
	engineerObject.Id = plan.Id.ValueString()
	engineerObject.Email = plan.Email.ValueString()

	// Print to standard logger; this will appear if TF_LOG=DEBUG is set
	log.Printf("Debug: Engineer Object: %#v", engineerObject)

	// Create new engineer
	engineer, err := r.client.CreateEngineer(engineerObject)
	if err != nil {
		log.Printf("Error: %v", err)
		resp.Diagnostics.AddError(
			"Error creating engineer",
			"Could not create engineer, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Name = types.StringValue(engineer.Name)
	plan.Id = types.StringValue(engineer.Id)
	plan.Email = types.StringValue(engineer.Email)

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
func (r *engineerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	log.Printf("Debug: Read request: %v", req)
	// Get current state
	var state engineerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		log.Printf("Error: %v", resp.Diagnostics)
		return
	}

	// Get refreshed engineer value from HashiCups
	engineer, err := r.client.GetEngineer(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error sending get request to devops-bootcamp api",
			"Could not read engineer Id "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state.Name = types.StringValue(engineer.Name)
	state.Id = types.StringValue(engineer.Id)
	state.Email = types.StringValue(engineer.Email)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *engineerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	log.Printf("Debug: Update request: %v", req)
	// Retrieve values from plan
	var plan engineerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var engineerObject devops_resource.Engineer
	engineerObject.Name = plan.Name.ValueString()
	engineerObject.Id = plan.Id.ValueString()
	engineerObject.Email = plan.Email.ValueString()

	// Print to standard logger; this will appear if TF_LOG=DEBUG is set
	log.Printf("Debug: Engineer Object: %#v", engineerObject)

	// Update existing engineer
	engineer, err := r.client.UpdateEngineer(engineerObject)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating engineer",
			"Could not update engineer, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Name = types.StringValue(engineer.Name)
	plan.Id = types.StringValue(engineer.Id)
	plan.Email = types.StringValue(engineer.Email)

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *engineerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state engineerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteEngineer(state.Id.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting engineer",
			"Could not delete order, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *engineerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *engineerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
