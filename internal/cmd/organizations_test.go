package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchOrganizations(t *testing.T) {
	tests := []struct {
		name      string
		mockData  interface{}
		mockError string
		wantOrgs  []Organization
		wantErr   bool
	}{
		{
			name: "正常系: 複数のOrgを返す",
			mockData: map[string]interface{}{
				"account": map[string]interface{}{
					"organizations": []map[string]string{
						{"id": "org1", "name": "Org One"},
						{"id": "org2", "name": "Org Two"},
					},
				},
			},
			wantOrgs: []Organization{
				{ID: "org1", Name: "Org One"},
				{ID: "org2", Name: "Org Two"},
			},
		},
		{
			name: "正常系: 空リスト",
			mockData: map[string]interface{}{
				"account": map[string]interface{}{
					"organizations": []map[string]string{},
				},
			},
			wantOrgs: []Organization{},
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

			orgs, err := fetchOrganizations(context.Background())
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantOrgs, orgs)
		})
	}
}

func TestResolveOrgID(t *testing.T) {
	defaultMockData := map[string]interface{}{
		"account": map[string]interface{}{
			"organizations": []map[string]string{
				{"id": "org1", "name": "Org One"},
				{"id": "org2", "name": "Org Two"},
			},
		},
	}
	emptyMockData := map[string]interface{}{
		"account": map[string]interface{}{
			"organizations": []map[string]string{},
		},
	}

	tests := []struct {
		name     string
		mockData interface{}
		orgFlag  string
		wantID   string
		wantName string
		wantErr  bool
	}{
		{
			name:     "デフォルト: --org未指定で先頭を返す",
			mockData: defaultMockData,
			orgFlag:  "",
			wantID:   "org1",
			wantName: "Org One",
		},
		{
			name:     "ID指定: IDでマッチ",
			mockData: defaultMockData,
			orgFlag:  "org2",
			wantID:   "org2",
			wantName: "Org Two",
		},
		{
			name:     "名前指定: 名前でマッチ",
			mockData: defaultMockData,
			orgFlag:  "Org One",
			wantID:   "org1",
			wantName: "Org One",
		},
		{
			name:     "不一致: マッチするorgなし",
			mockData: defaultMockData,
			orgFlag:  "nonexistent",
			wantErr:  true,
		},
		{
			name:     "空リスト: orgが存在しない",
			mockData: emptyMockData,
			orgFlag:  "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := newMockGraphQLServer(t, tt.mockData)
			defer srv.Close()
			setupTest(t, srv.URL)
			orgFlag = tt.orgFlag

			id, name, err := resolveOrgID(context.Background())
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantID, id)
			assert.Equal(t, tt.wantName, name)
		})
	}
}
