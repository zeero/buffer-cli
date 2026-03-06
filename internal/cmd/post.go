package cmd

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeero/buffer-v2-cli/client"
)

type PostResult struct {
	ChannelID   string `json:"channelId"`
	ChannelName string `json:"channelName"`
	PostID      string `json:"postId,omitempty"`
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
}

var postText string

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "全チャンネルにテキストを投稿する",
	Run: func(cmd *cobra.Command, args []string) {
		if apiToken == "" {
			printError("APIトークンが設定されていません。--token または BUFFER_TOKEN 環境変数を設定してください")
		}
		if postText == "" {
			printError("--text は必須です")
		}

		ctx := context.Background()
		orgID, _, err := resolveOrgID(ctx)
		if err != nil {
			printError(err.Error())
		}

		channels, err := fetchChannels(ctx, orgID)
		if err != nil {
			printError(err.Error())
		}
		if len(channels) == 0 {
			printError("投稿先チャンネルが見つかりません")
		}

		results := make([]PostResult, 0, len(channels))
		for _, ch := range channels {
			postID, postErr := createPost(ctx, ch.ID, postText)
			r := PostResult{
				ChannelID:   ch.ID,
				ChannelName: ch.Name,
				Success:     postErr == nil,
			}
			if postErr == nil {
				r.PostID = postID
			} else {
				r.Error = postErr.Error()
			}
			results = append(results, r)
		}

		printSuccess(func() {
			color.Cyan("=== Posting to %d channels ===", len(channels))
			for _, r := range results {
				if r.Success {
					fmt.Printf("  %s %s  → %s\n", color.GreenString("✓"), r.ChannelName, r.PostID)
				} else {
					fmt.Printf("  %s %s  → failed: %s\n", color.RedString("✗"), r.ChannelName, r.Error)
				}
			}
		}, results)
	},
}

func createPost(ctx context.Context, channelID, text string) (string, error) {
	c := client.NewWithEndpoint(apiToken, bufferEndpoint)

	// schedulingType: automatic, mode: now で即時投稿
	// レスポンスは PostActionSuccess | MutationError のunion型
	mutation := `
	mutation CreatePost($channelId: ChannelId!, $text: String!) {
		createPost(input: {
			channelId: $channelId
			text: $text
			schedulingType: automatic
			mode: now
		}) {
			... on PostActionSuccess {
				post {
					id
				}
			}
			... on MutationError {
				message
			}
		}
	}`

	var resp struct {
		CreatePost struct {
			Post    struct{ ID string `json:"id"` } `json:"post"`
			Message string                          `json:"message"`
		} `json:"createPost"`
	}

	vars := map[string]interface{}{
		"channelId": channelID,
		"text":      text,
	}

	if err := c.Run(ctx, mutation, vars, &resp); err != nil {
		return "", err
	}
	if resp.CreatePost.Message != "" {
		return "", fmt.Errorf("%s", resp.CreatePost.Message)
	}

	return resp.CreatePost.Post.ID, nil
}

func init() {
	postCmd.Flags().StringVar(&postText, "text", "", "投稿するテキスト（必須）")
	rootCmd.AddCommand(postCmd)
}
