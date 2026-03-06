package cmd

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeero/buffer-v2-cli/client"
)

type Channel struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Service string `json:"service"`
}

var channelsCmd = &cobra.Command{
	Use:   "channels",
	Short: "チャンネル一覧を取得する",
	Run: func(cmd *cobra.Command, args []string) {
		if apiToken == "" {
			printError("APIトークンが設定されていません。--token または BUFFER_TOKEN 環境変数を設定してください")
		}

		ctx := context.Background()
		orgID, orgName, err := resolveOrgID(ctx)
		if err != nil {
			printError(err.Error())
		}

		channels, err := fetchChannels(ctx, orgID)
		if err != nil {
			printError(err.Error())
		}

		printSuccess(func() {
			color.Cyan("=== Channels (%s) ===", orgName)
			for _, ch := range channels {
				fmt.Printf("  • %s [%s]  (%s)\n", color.GreenString(ch.Name), ch.Service, ch.ID)
			}
		}, channels)
	},
}

func fetchChannels(ctx context.Context, organizationID string) ([]Channel, error) {
	c := client.NewWithEndpoint(apiToken, bufferEndpoint)

	query := `
	query GetChannels($organizationId: OrganizationId!) {
		channels(input: {
			organizationId: $organizationId
		}) {
			id
			name
			service
		}
	}`

	var resp struct {
		Channels []Channel `json:"channels"`
	}

	vars := map[string]interface{}{
		"organizationId": organizationID,
	}

	if err := c.Run(ctx, query, vars, &resp); err != nil {
		return nil, err
	}

	return resp.Channels, nil
}

func init() {
	rootCmd.AddCommand(channelsCmd)
}
