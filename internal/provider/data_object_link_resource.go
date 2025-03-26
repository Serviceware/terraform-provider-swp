package provider

import (
	"context"
	"fmt"

	"github.com/Serviceware/terraform-provider-swp/internal/aipe"
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
		MarkdownDescription: "Creates a link between two data objects",
		Attributes: map[string]schema.Attribute{
			"source_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The system.id of the source object",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"link_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the link. This is the name of the link ('incident-causes'), not the name of the relation (of which a link has 2 - 'caused-by' or 'causes').",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"relation_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the relation. This is the 'end' of the link on the source objects side",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "This is the list of target object IDs to link to",
				Required:            true,
			},
		},
	}
}

func (r *DataObjectLinkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	resp.TypeName = req.ProviderTypeName + "_aipe_data_object_link"
}

func (d *DataObjectLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DataObjectLinkResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Target IDs", map[string]interface{}{"target_ids": data.TargetIDs})
	if len(data.TargetIDs) > 0 {
		// Create the link in the AIPE API.
		tflog.Info(ctx, "Creating link", map[string]interface{}{"data": data.SourceID.ValueString()})
		err := d.client.UpdateDataObjectLinks(ctx, data.SourceID.ValueString(), data.LinkName.ValueString(), data.RelationName.ValueString(), data.TargetIDs, nil)
		if err != nil {
			resp.Diagnostics.AddError("Failed to create link", err.Error())
		}
	} else {
		tflog.Info(ctx, "No link to create", map[string]interface{}{"data": data.SourceID.ValueString()})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataObjectLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DataObjectLinkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	linkData, err := d.client.GetDataObjectLinks(ctx, data.SourceID.ValueString(), data.LinkName.ValueString(), data.RelationName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get link data", err.Error())
		return
	}

	if linkData == nil {
		data.TargetIDs = []string{}
	} else {
		data.TargetIDs = linkData
	}

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

	var stateTargetIDs = state.TargetIDs
	var planTargetIDs = plan.TargetIDs

	add, remove := diffStateAndPlanIDs(stateTargetIDs, planTargetIDs)

	if len(add) == 0 && len(remove) == 0 {
		tflog.Info(ctx, "No changes to link", map[string]interface{}{"data": plan.SourceID.ValueString()})
	} else {
		tflog.Info(ctx, "Updating link", map[string]interface{}{"add": add, "remove": remove})
		err := d.client.UpdateDataObjectLinks(ctx, plan.SourceID.ValueString(), plan.LinkName.ValueString(), plan.RelationName.ValueString(), add, remove)
		if err != nil {
			resp.Diagnostics.AddError("Failed to update link", err.Error())
		}
	}
	state.TargetIDs = plan.TargetIDs

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func diffStateAndPlanIDs(stateTargetIDs []string, planTargetIDs []string) ([]string, []string) {
	var allIDs = make(map[string]bool)
	var targetIDsInPlan = make(map[string]bool)
	var targetIDsInState = make(map[string]bool)

	for _, id := range stateTargetIDs {
		allIDs[id] = true
		targetIDsInState[id] = true
	}
	for _, id := range planTargetIDs {
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
	return add, remove
}

func (d *DataObjectLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DataObjectLinkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(data.TargetIDs) > 0 {
		tflog.Info(ctx, "Deleting link", map[string]interface{}{"data": data.SourceID.ValueString()})
		err := d.client.UpdateDataObjectLinks(ctx, data.SourceID.ValueString(), data.LinkName.ValueString(), data.RelationName.ValueString(), nil, data.TargetIDs)
		if err != nil {
			resp.Diagnostics.AddError("Failed to delete link", err.Error())
		}
	} else {
		tflog.Info(ctx, "No link to delete", map[string]interface{}{"data": data.SourceID.ValueString()})
	}
}
