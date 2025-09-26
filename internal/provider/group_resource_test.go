package provider

import (
	"context"
	"testing"

	"github.com/extenda/terraform-provider-hiiretail-iam/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// IClient is an interface for the client.Client type
type IClient interface {
	CreateGroup(ctx context.Context, name string, description string) (*client.Group, error)
	GetGroup(ctx context.Context, id string) (*client.Group, error)
	UpdateGroup(ctx context.Context, id string, name string, description string) (*client.Group, error)
	DeleteGroup(ctx context.Context, id string) error
}

// MockClient is a mock implementation of the IClient interface
type MockClient struct {
	mock.Mock
}

func NewMockClient() client.Client {
	c := client.NewClient("https://test-api.com", "test-token")
	return *c
}

func (m *MockClient) CreateGroup(ctx context.Context, name string, description string) (*client.Group, error) {
	args := m.Called(ctx, name, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.Group), args.Error(1)
}

func (m *MockClient) GetGroup(ctx context.Context, id string) (*client.Group, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.Group), args.Error(1)
}

func (m *MockClient) UpdateGroup(ctx context.Context, id string, name string, description string) (*client.Group, error) {
	args := m.Called(ctx, id, name, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.Group), args.Error(1)
}

func (m *MockClient) DeleteGroup(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestNewGroupResource(t *testing.T) {
	r := NewGroupResource()
	assert.NotNil(t, r)
}

func TestGroupResourceSchema(t *testing.T) {
	r := &GroupResource{}
	response := &resource.SchemaResponse{}

	r.Schema(context.Background(), resource.SchemaRequest{}, response)

	assert.NotNil(t, response.Schema)
	assert.NotNil(t, response.Schema.Attributes)

	// Check id attribute
	assert.Contains(t, response.Schema.Attributes, "id")
	idAttr := response.Schema.Attributes["id"].(schema.StringAttribute)
	assert.True(t, idAttr.Computed)
	assert.False(t, idAttr.Required)
	assert.False(t, idAttr.Optional)

	// Check name attribute
	assert.Contains(t, response.Schema.Attributes, "name")
	nameAttr := response.Schema.Attributes["name"].(schema.StringAttribute)
	assert.True(t, nameAttr.Required)
	assert.False(t, nameAttr.Computed)
	assert.False(t, nameAttr.Optional)

	// Check description attribute
	assert.Contains(t, response.Schema.Attributes, "description")
	descAttr := response.Schema.Attributes["description"].(schema.StringAttribute)
	assert.True(t, descAttr.Required)
	assert.False(t, descAttr.Computed)
	assert.False(t, descAttr.Optional)
}

func TestGroupResource_Create(t *testing.T) {
	mockClient := &MockClient{}
	r := &GroupResource{}

	expectedGroup := &client.Group{
		ID:          "test-id",
		Name:        "test-group",
		Description: "test description",
	}

	// Set up mock expectations before configuring the resource
	mockClient.On("CreateGroup", mock.Anything, "test-group", "test description").Return(expectedGroup, nil)

	// Configure the resource with the mock client
	configResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: mockClient,
	}, configResp)
	assert.False(t, configResp.Diagnostics.HasError())

	ctx := context.Background()

	// Get schema from resource
	schemaResp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, schemaResp)

	// Test with valid configuration
	plan := tfsdk.Plan{
		Schema: schemaResp.Schema,
	}
	_ = plan.Set(ctx, &GroupResourceModel{
		Name:        types.StringValue("test-group"),
		Description: types.StringValue("test description"),
	})

	req := resource.CreateRequest{
		Plan: plan,
	}

	state := tfsdk.State{
		Schema: schemaResp.Schema,
	}
	createResp := &resource.CreateResponse{
		State: state,
	}

	r.Create(ctx, req, createResp)

	// Verify
	assert.False(t, createResp.Diagnostics.HasError())
	mockClient.AssertExpectations(t)

	// Check the state
	var actualState GroupResourceModel
	_ = createResp.State.Get(ctx, &actualState)
	assert.Equal(t, "test-id", actualState.ID.ValueString())
	assert.Equal(t, "test-group", actualState.Name.ValueString())
	assert.Equal(t, "test description", actualState.Description.ValueString())
}

func TestGroupResource_Read(t *testing.T) {
	mockClient := &MockClient{}
	r := &GroupResource{}

	configResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: mockClient,
	}, configResp)
	assert.False(t, configResp.Diagnostics.HasError())

	expectedGroup := &client.Group{
		ID:          "test-id",
		Name:        "test-group",
		Description: "test description",
	}

	mockClient.On("GetGroup", mock.Anything, "test-id").Return(expectedGroup, nil)

	ctx := context.Background()

	// Get schema from resource
	schemaResp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, schemaResp)

	// Test with valid state
	state := tfsdk.State{
		Schema: schemaResp.Schema,
	}
	_ = state.Set(ctx, &GroupResourceModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("old-name"),
		Description: types.StringValue("old description"),
	})

	req := resource.ReadRequest{
		State: state,
	}
	resp := &resource.ReadResponse{
		State: state,
	}

	r.Read(ctx, req, resp)

	// Verify
	assert.False(t, resp.Diagnostics.HasError())
	mockClient.AssertExpectations(t)

	// Check the state
	var actualState GroupResourceModel
	_ = resp.State.Get(ctx, &actualState)
	assert.Equal(t, "test-id", actualState.ID.ValueString())
	assert.Equal(t, "test-group", actualState.Name.ValueString())
	assert.Equal(t, "test description", actualState.Description.ValueString())
}