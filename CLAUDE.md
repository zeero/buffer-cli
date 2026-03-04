# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## コマンド

```bash
# ビルド
go build ./...

# 実行（開発時）
go run . <command> [flags]

# 依存関係追加
go get <package>
go mod tidy
```

## アーキテクチャ

Buffer API（GraphQL）を操作する CLI ツール。認証トークンは `--token` フラグまたは `$BUFFER_TOKEN` 環境変数で渡す。

### データフロー

`post` コマンドは次の順序で API を呼び出す:

```
GetOrganizations → resolveOrgID → GetChannels → CreatePost (× チャンネル数)
```

`channels` コマンドも同様に `resolveOrgID` を経由する。

### 共有関数（`cmd/` パッケージ内）

- `resolveOrgID(ctx)` — `--org` フラグの値（ID または名前）で org を解決。省略時は先頭 1 件を返す。`organizations.go` に定義。
- `fetchOrganizations(ctx)` / `fetchChannels(ctx, orgID)` — 各コマンドファイルに定義し、複数コマンドから呼び出される。
- `printSuccess(humanOutput, data)` / `printError(msg)` — `--json` フラグに応じて出力形式を切り替える。`output.go` に定義。

### 出力形式

`--json` フラグで人間向けとエージェント向けを切り替える。JSON は常に `{"success": bool, "data": ..., "error": "..."}` の統一フォーマット。

### GraphQL クライアント

`client/buffer.go` の `BufferClient.Run()` がすべての API 呼び出しを担う。`Authorization: Bearer <token>` ヘッダーを自動付与。API がベータ版のため `machinebox/graphql`（動的クライアント）を採用。安定後は `Khan/genqlient` への移行を検討。

### 新コマンド追加時のパターン

1. `cmd/<name>.go` を作成
2. `var <name>Cmd = &cobra.Command{...}` を定義
3. `init()` で `rootCmd.AddCommand(<name>Cmd)` を呼ぶ
4. `apiToken` / `jsonOutput` / `orgFlag` はパッケージ変数として `root.go` で管理されているため、そのまま参照可能
