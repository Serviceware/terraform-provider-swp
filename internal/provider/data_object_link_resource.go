package provider

import (
	"context"
	"fmt"

	"github.com/Serviceware/terraform-provider-aipe/internal/aipe"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &DataObjectLinkResource{}

func NewDataObjectLinkResource() resource.Resource {
	return &DataObjectLinkResource{}
}

type DataObjectLinkResource struct {
	client *aipe.AIPEClient
}

type DataObjectLinkResourceModel struct {
	SourceID types.String `tfsdk:"source_id"`

	// The overall link name, e.g. "Requester", "Incident Cause", ...
	LinkName types.String `tfsdk:"link_name"`

	// The name of the "side" of the link, e.g. "requested by" or "causes
	RelationName types.String `tfsdk:"relation_name"`

	TargetIDs []string `tfsdk:"target_ids"`
}

func (d *DataObjectLinkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Link between two data objects",
		Attributes: map[string]schema.Attribute{
			"source_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"link_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"relation_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

func (r *DataObjectLinkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*aipe.AIPEClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *aipe.AIPEClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (d *DataObjectLinkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_object_link"
}

func (d *DataObjectLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DataObjectLinkResourceModel

	// Read the data from the request.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the link in the AIPE API.
	tflog.Info(ctx, "Creating link", map[string]interface{}{"data": data.SourceID.ValueString()})
	err := d.client.UpdateDataObjectLinks(ctx, data.SourceID.ValueString(), data.LinkName.ValueString(), data.RelationName.ValueString(), data.TargetIDs, nil)

	if err != nil {
		resp.Diagnostics.AddError("Failed to create link", err.Error())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataObjectLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DataObjectLinkResourceModel

	// Read the data from the AIPE API.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the link data from the AIPE API.
	linkData, err := d.client.GetDataObjectLinks(ctx, data.SourceID.ValueString(), data.LinkName.ValueString(), data.RelationName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get link data", err.Error())
		return
	}

	data.TargetIDs = linkData

	// Write the data to the response.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataObjectLinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DataObjectLinkResourceModel
	var state DataObjectLinkResourceModel

	// Read the data from the request.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var allIDs = make(map[string]bool)
	var targetIDsInPlan = make(map[string]bool)
	var targetIDsInState = make(map[string]bool)

	for _, id := range state.TargetIDs {
		allIDs[id] = true
		targetIDsInState[id] = true
	}
	for _, id := range plan.TargetIDs {
		allIDs[id] = true
		targetIDsInPlan[id] = true
	}

	var add []string
	var remove []string

	for id := range allIDs {
		if targetIDsInPlan[id] && !targetIDsInState[id] {
			add = append(add, id)
		}

		if !targetIDsInPlan[id] && targetIDsInState[id] {
			remove = append(remove, id)
		}
	}

	// Update the link in the AIPE API.
	tflog.Info(ctx, "Updating link", map[string]interface{}{"add": add, "remove": remove})
	err := d.client.UpdateDataObjectLinks(ctx, plan.SourceID.ValueString(), plan.LinkName.ValueString(), plan.RelationName.ValueString(), add, remove)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update link", err.Error())
	}
	state.TargetIDs = plan.TargetIDs

	// Write the data to the response.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *DataObjectLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DataObjectLinkResourceModel

	// Read the data from the request.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the link in the AIPE API.
	tflog.Info(ctx, "Deleting link", map[string]interface{}{"data": data.SourceID.ValueString()})
	err := d.client.UpdateDataObjectLinks(ctx, data.SourceID.ValueString(), data.LinkName.ValueString(), data.RelationName.ValueString(), nil, data.TargetIDs)

	if err != nil {
		resp.Diagnostics.AddError("Failed to delete link", err.Error())
	}
	panic("unimplemented")
}
