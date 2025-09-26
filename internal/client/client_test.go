package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CreateGroup(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/groups", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": "test-id",
			"name": "test-group",
			"description": "test description"
		}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-token")
	group, err := client.CreateGroup(context.Background(), "test-group", "test description")

	assert.NoError(t, err)
	assert.Equal(t, "test-id", group.ID)
	assert.Equal(t, "test-group", group.Name)
	assert.Equal(t, "test description", group.Description)
}

func TestClient_GetGroup(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/groups/test-id", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "test-id",
			"name": "test-group",
			"description": "test description"
		}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-token")
	group, err := client.GetGroup(context.Background(), "test-id")

	assert.NoError(t, err)
	assert.Equal(t, "test-id", group.ID)
	assert.Equal(t, "test-group", group.Name)
	assert.Equal(t, "test description", group.Description)
}

func TestClient_UpdateGroup(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/groups/test-id", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "test-id",
			"name": "updated-group",
			"description": "updated description"
		}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-token")
	group, err := client.UpdateGroup(context.Background(), "test-id", "updated-group", "updated description")

	assert.NoError(t, err)
	assert.Equal(t, "test-id", group.ID)
	assert.Equal(t, "updated-group", group.Name)
	assert.Equal(t, "updated description", group.Description)
}

func TestClient_DeleteGroup(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/groups/test-id", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-token")
	err := client.DeleteGroup(context.Background(), "test-id")

	assert.NoError(t, err)
}