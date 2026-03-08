package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/beezly/terraform-provider-unifi/internal/client"
)

// Ensure interface compliance.
var _ resource.Resource = &NetworkResource{}
var _ resource.ResourceWithImportState = &NetworkResource{}

// NetworkResource implements the unifi_network resource.
type NetworkResource struct {
	client *client.Client
}

func NewNetworkResource() resource.Resource {
	return &NetworkResource{}
}

// NetworkModel maps to the UniFi network API response.
type NetworkModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	VlanID      types.Int64  `tfsdk:"vlan_id"`
	IPSubnet    types.String `tfsdk:"ip_subnet"`
	DhcpEnabled types.Bool   `tfsdk:"dhcp_enabled"`
}

// networkAPIPayload is what we send/receive from the UniFi API.
type networkAPIPayload struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Enabled     *bool  `json:"enabled,omitempty"`
	VlanID      *int64 `json:"vlanId,omitempty"`
	IPSubnet    string `json:"ipSubnet,omitempty"`
	DHCPEnabled *bool  `json:"dhcpEnabled,omitempty"`
}

type networkAPIResponse struct {
	Data networkAPIPayload `json:"data"`
}

func (r *NetworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

func (r *NetworkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a UniFi network (VLAN).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Network ID (UUID).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Network name.",
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether the network is enabled.",
			},
			"vlan_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "VLAN ID (1–4094). Omit for the default untagged network.",
			},
			"ip_subnet": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "IPv4 subnet in CIDR notation (e.g. `192.168.3.0/24`).",
			},
			"dhcp_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable DHCP server on this network.",
			},
		},
	}
}

func (r *NetworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data type", fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *NetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := modelToPayload(plan)
	var result networkAPIResponse
	err := r.client.Post(ctx,
		fmt.Sprintf("/v1/sites/%s/networks", r.client.SiteID()),
		payload, &result)
	if err != nil {
		resp.Diagnostics.AddError("Error creating network", err.Error())
		return
	}

	payloadToModel(&result.Data, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result networkAPIResponse
	err := r.client.Get(ctx,
		fmt.Sprintf("/v1/sites/%s/networks/%s", r.client.SiteID(), state.ID.ValueString()),
		&result)
	if err != nil {
		resp.Diagnostics.AddError("Error reading network", err.Error())
		return
	}

	payloadToModel(&result.Data, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := modelToPayload(plan)
	var result networkAPIResponse
	err := r.client.Put(ctx,
		fmt.Sprintf("/v1/sites/%s/networks/%s", r.client.SiteID(), plan.ID.ValueString()),
		payload, &result)
	if err != nil {
		resp.Diagnostics.AddError("Error updating network", err.Error())
		return
	}

	payloadToModel(&result.Data, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NetworkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx,
		fmt.Sprintf("/v1/sites/%s/networks/%s", r.client.SiteID(), state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting network", err.Error())
	}
}

func (r *NetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by network UUID
	var state NetworkModel
	state.ID = types.StringValue(req.ID)

	var result networkAPIResponse
	err := r.client.Get(ctx,
		fmt.Sprintf("/v1/sites/%s/networks/%s", r.client.SiteID(), req.ID),
		&result)
	if err != nil {
		resp.Diagnostics.AddError("Error importing network", err.Error())
		return
	}

	payloadToModel(&result.Data, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// modelToPayload converts the Terraform model to the API payload.
func modelToPayload(m NetworkModel) networkAPIPayload {
	p := networkAPIPayload{
		Name: m.Name.ValueString(),
	}
	if !m.Enabled.IsNull() && !m.Enabled.IsUnknown() {
		v := m.Enabled.ValueBool()
		p.Enabled = &v
	}
	if !m.VlanID.IsNull() && !m.VlanID.IsUnknown() {
		v := m.VlanID.ValueInt64()
		p.VlanID = &v
	}
	if !m.IPSubnet.IsNull() && !m.IPSubnet.IsUnknown() {
		p.IPSubnet = m.IPSubnet.ValueString()
	}
	if !m.DhcpEnabled.IsNull() && !m.DhcpEnabled.IsUnknown() {
		v := m.DhcpEnabled.ValueBool()
		p.DHCPEnabled = &v
	}
	return p
}

// payloadToModel maps the API response back into the Terraform model.
func payloadToModel(p *networkAPIPayload, m *NetworkModel) {
	m.ID = types.StringValue(p.ID)
	m.Name = types.StringValue(p.Name)
	if p.Enabled != nil {
		m.Enabled = types.BoolValue(*p.Enabled)
	}
	if p.VlanID != nil {
		m.VlanID = types.Int64Value(*p.VlanID)
	}
	if p.IPSubnet != "" {
		m.IPSubnet = types.StringValue(p.IPSubnet)
	}
	if p.DHCPEnabled != nil {
		m.DhcpEnabled = types.BoolValue(*p.DHCPEnabled)
	}
}
