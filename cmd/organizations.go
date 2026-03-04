package cmd

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeero/buffer-v2-cli/client"
)

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var organizationsCmd = &cobra.Command{
	Use:   "organizations",
	Short: "オーガニゼーション一覧を取得する",
	Run: func(cmd *cobra.Command, args []string) {
		if apiToken == "" {
			printError("APIトークンが設定されていません。--token または BUFFER_TOKEN 環境変数を設定してください")
		}

		orgs, err := fetchOrganizations(context.Background())
		if err != nil {
			printError(err.Error())
		}

		printSuccess(func() {
			color.Cyan("=== Organizations ===")
			for _, org := range orgs {
				fmt.Printf("  • %s  (%s)\n", color.GreenString(org.Name), org.ID)
			}
		}, orgs)
	},
}

func fetchOrganizations(ctx context.Context) ([]Organization, error) {
	c := client.NewWithEndpoint(apiToken, bufferEndpoint)

	query := `
	query GetOrganizations {
		account {
			organizations {
				id
				name
			}
		}
	}`

	var resp struct {
		Account struct {
			Organizations []Organization `json:"organizations"`
		} `json:"account"`
	}

	if err := c.Run(ctx, query, nil, &resp); err != nil {
		return nil, err
	}

	return resp.Account.Organizations, nil
}

// resolveOrgID は --org フラグの値または先頭 org の ID を返す
func resolveOrgID(ctx context.Context) (string, string, error) {
	orgs, err := fetchOrganizations(ctx)
	if err != nil {
		return "", "", err
	}
	if len(orgs) == 0 {
		return "", "", fmt.Errorf("オーガニゼーションが見つかりません")
	}

	if orgFlag == "" {
		return orgs[0].ID, orgs[0].Name, nil
	}

	for _, org := range orgs {
		if org.ID == orgFlag || org.Name == orgFlag {
			return org.ID, org.Name, nil
		}
	}

	return "", "", fmt.Errorf("オーガニゼーション '%s' が見つかりません", orgFlag)
}

func init() {
	rootCmd.AddCommand(organizationsCmd)
}
