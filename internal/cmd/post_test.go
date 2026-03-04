package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name      string
		mockData  interface{}
		mockError string
		wantID    string
		wantErr   bool
	}{
		{
			name: "正常系: PostIDを返す",
			mockData: map[string]interface{}{
				"createPost": map[string]interface{}{
					"post": map[string]string{"id": "post123"},
				},
			},
			wantID: "post123",
		},
		{
			name:      "エラー系: GraphQLエラー",
			mockError: "channel not found",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := newMockServer(t, tt.mockData, tt.mockError)
			defer srv.Close()
			setupTest(t, srv.URL)

			postID, err := createPost(context.Background(), "ch1", "test text")
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantID, postID)
		})
	}
}

func TestCreatePost_RequestPayload(t *testing.T) {
	var gotBody map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"createPost":{"post":{"id":"post456"}}}}`))
	}))
	defer srv.Close()
	setupTest(t, srv.URL)

	_, err := createPost(context.Background(), "ch42", "hello world")
	require.NoError(t, err)

	vars, ok := gotBody["variables"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "ch42", vars["channelId"])
	assert.Equal(t, "hello world", vars["text"])
}
