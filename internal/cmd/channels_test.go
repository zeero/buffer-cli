package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchChannels(t *testing.T) {
	tests := []struct {
		name         string
		mockData     interface{}
		mockError    string
		wantChannels []Channel
		wantErr      bool
	}{
		{
			name: "正常系: チャンネル一覧を返す",
			mockData: map[string]interface{}{
				"channels": []map[string]string{
					{"id": "ch1", "name": "Twitter", "service": "twitter"},
					{"id": "ch2", "name": "Facebook", "service": "facebook"},
				},
			},
			wantChannels: []Channel{
				{ID: "ch1", Name: "Twitter", Service: "twitter"},
				{ID: "ch2", Name: "Facebook", Service: "facebook"},
			},
		},
		{
			name: "正常系: 空リスト",
			mockData: map[string]interface{}{
				"channels": []map[string]string{},
			},
			wantChannels: []Channel{},
		},
		{
			name:      "エラー系: GraphQLエラー",
			mockError: "unauthorized",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := newMockServer(t, tt.mockData, tt.mockError)
			defer srv.Close()
			setupTest(t, srv.URL)

			channels, err := fetchChannels(context.Background(), "org1")
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantChannels, channels)
		})
	}
}
