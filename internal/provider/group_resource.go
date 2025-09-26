package provider

import (
	"context"
	"fmt"

	"github.com/extenda/terraform-provider-hiiretail-iam/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &GroupResource{}
var _ resource.ResourceWithImportState = &GroupResource{}

func NewGroupResource() resource.Resource {
	return &GroupResource{
		logger: client.NewLogger(client.LogLevelInfo),
	}
}

// GroupResource defines the implementation of the hiiretail_iam_group resource.
// It handles the lifecycle (create, read, update, delete) of IAM groups
// in the HiiRetail system, ensuring proper state management and error handling.
type GroupResource struct {
	client client.IClient
	logger *client.Logger
}

// GroupResourceModel describes the resource data model for an IAM group.
// It maps the API response to Terraform's schema attributes, handling
// both required and computed fields in a type-safe manner using the
// Terraform Plugin Framework's types.
type GroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (r *GroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *GroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IAM group in HiiRetail. Groups are collections of users that share the same access permissions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The group identifier.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the group. Must be between 3 and 128 characters, start with a letter, and contain only letters, numbers, hyphens (-), and underscores (_).",
				Required:    true,
				Validators: []validator.String{
					groupNameValidator{},
				},
			},
			"description": schema.StringAttribute{
				Description: "A description of the group explaining its purpose. Must be between 1 and 256 characters.",
				Required:    true,
				Validators: []validator.String{
					groupDescriptionValidator{},
				},
			},
		},
	}
}

func (r *GroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(client.IClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected client.IClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new group
	group, err := r.client.CreateGroup(ctx, data.Name.ValueString(), data.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to Create Group", fmt.Sprintf("Could not create IAM group '%s'. This might be due to a name conflict or invalid input. Original error: %s", data.Name.ValueString(), err))
		return
	}

	// Map response body to schema
	data.ID = types.StringValue(group.ID)
	data.Name = types.StringValue(group.Name)
	data.Description = types.StringValue(group.Description)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed group value from HiiRetail
	group, err := r.client.GetGroup(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsResourceNotFound(err) {
			// If the resource does not exist, remove it from state
			r.logger.Info("Group %s no longer exists", data.ID.ValueString())
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to Read Group", fmt.Sprintf("Could not read IAM group (ID: %s). This might be due to insufficient permissions or network issues. Original error: %s", data.ID.ValueString(), err))
		return
	}

	// Map response body to schema
	data.ID = types.StringValue(group.ID)
	data.Name = types.StringValue(group.Name)
	data.Description = types.StringValue(group.Description)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing group
	group, err := r.client.UpdateGroup(ctx, data.ID.ValueString(), data.Name.ValueString(), data.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to Update Group", fmt.Sprintf("Could not update IAM group '%s' (ID: %s). This might be due to concurrent modifications or invalid input. Original error: %s", data.Name.ValueString(), data.ID.ValueString(), err))
		return
	}

	// Map response body to schema
	data.Name = types.StringValue(group.Name)
	data.Description = types.StringValue(group.Description)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing group
	err := r.client.DeleteGroup(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsResourceNotFound(err) {
			// If the resource is already gone, that's okay
			r.logger.Info("Group %s already deleted", data.ID.ValueString())
			return
		}
		resp.Diagnostics.AddError("Failed to Delete Group", fmt.Sprintf("Could not delete IAM group (ID: %s). This might be due to the group having existing members or dependencies. Original error: %s", data.ID.ValueString(), err))
		return
	}
}

func (r *GroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Read the group data directly
	group, err := r.client.GetGroup(ctx, req.ID)
	if err != nil {
		if client.IsResourceNotFound(err) {
			resp.Diagnostics.AddError("Resource Not Found", 
				fmt.Sprintf("Unable to find group %s for import", req.ID))
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read group during import, got error: %s", err))
		return
	}

	// Set into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &GroupResourceModel{
		ID:          types.StringValue(group.ID),
		Name:        types.StringValue(group.Name),
		Description: types.StringValue(group.Description),
	})...)
}