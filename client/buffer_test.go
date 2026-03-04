package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_SendsAuthHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{}}`))
	}))
	defer srv.Close()

	c := NewWithEndpoint("test-token", srv.URL)
	var resp struct{}
	err := c.Run(context.Background(), `query {}`, nil, &resp)
	require.NoError(t, err)
	assert.Equal(t, "Bearer test-token", gotAuth)
}

func TestRun_DeserializesResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"name":"test-value"}}`))
	}))
	defer srv.Close()

	c := NewWithEndpoint("token", srv.URL)
	var resp struct {
		Name string `json:"name"`
	}
	err := c.Run(context.Background(), `query {}`, nil, &resp)
	require.NoError(t, err)
	assert.Equal(t, "test-value", resp.Name)
}

func TestRun_ReturnsErrorOnGraphQLError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"errors":[{"message":"something went wrong"}]}`))
	}))
	defer srv.Close()

	c := NewWithEndpoint("token", srv.URL)
	var resp struct{}
	err := c.Run(context.Background(), `query {}`, nil, &resp)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "something went wrong")
}

func TestRun_SendsVariables(t *testing.T) {
	var gotBody map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{}}`))
	}))
	defer srv.Close()

	c := NewWithEndpoint("token", srv.URL)
	vars := map[string]interface{}{"key": "value"}
	var resp struct{}
	err := c.Run(context.Background(), `query {}`, vars, &resp)
	require.NoError(t, err)

	variables, ok := gotBody["variables"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "value", variables["key"])
}
