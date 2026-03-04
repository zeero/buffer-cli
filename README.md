Buffer の API をターミナルから手軽に操作できる CLI ツールです。
組織・チャンネルの確認から、複数チャンネルへの一括投稿まで、1コマンドで実行できます。

## 目次

1. [インストール](#インストール)
   1. [go install（推奨）](#go-install推奨)
   1. [ソースからビルド](#ソースからビルド)
1. [セットアップ](#セットアップ)
1. [使い方](#使い方)
1. [コマンドリファレンス](#コマンドリファレンス)
   1. [共通フラグ](#共通フラグ)
   1. [organizations](#organizations)
   1. [channels](#channels)
   1. [post](#post)
1. [開発](#開発)

---

## インストール

### go install（推奨）

```bash
go install github.com/zeero/buffer-v2-cli@latest
```

インストール後、`buffer` コマンドが使えるようになります。
`$GOPATH/bin`（通常 `~/go/bin`）を `$PATH` に追加しておいてください。

### ソースからビルド

```bash
git clone https://github.com/zeero/buffer-v2-cli.git
cd buffer-v2-cli
go build -o buffer .
```

## セットアップ

Buffer の API トークンが必要です。[Buffer の設定画面](https://buffer.com) から取得してください。

トークンは **環境変数** または **フラグ** で渡せます。

```bash
# 環境変数（推奨）
export BUFFER_TOKEN=your_token_here

# フラグで都度指定する場合
buffer --token your_token_here <command>
```

> [!TIP]
> `BUFFER_TOKEN` を `.zshrc` / `.bashrc` に追加しておくと、毎回 `--token` を省略できます。

## 使い方

```bash
# 所属オーガニゼーションを確認
buffer organizations

# チャンネル一覧を確認
buffer channels

# 全チャンネルにテキストを投稿
buffer post --text "Hello from CLI!"

# 特定のオーガニゼーションを指定して投稿
buffer post --org "My Org" --text "Hello!"

# JSON 形式で出力（スクリプト・AI エージェント連携向け）
buffer channels --json | jq '.data[].name'
```

## コマンドリファレンス

### 共通フラグ

すべてのコマンドで使用できます。

| フラグ | 環境変数 | デフォルト | 説明 |
|---|---|---|---|
| `--token` | `BUFFER_TOKEN` | — | Buffer API トークン（必須） |
| `--org` | — | 先頭の1件 | オーガニゼーション ID または名前 |
| `--json` | — | `false` | JSON 形式で出力する |

### organizations

所属しているオーガニゼーションの一覧を取得します。

```bash
buffer organizations [--json]
```

```
=== Organizations ===
  • My Company  (org_xxxxxxxx)
  • Personal    (org_yyyyyyyy)
```

### channels

オーガニゼーションに紐づくチャンネル（SNS アカウント）の一覧を取得します。

```bash
buffer channels [--org <id|name>] [--json]
```

```
=== Channels (My Company) ===
  • @myaccount [twitter]  (ch_xxxxxxxx)
  • My Page    [facebook] (ch_yyyyyyyy)
```

### post

オーガニゼーションの**全チャンネル**に対してテキストを即時投稿します。

```bash
buffer post --text <text> [--org <id|name>] [--json]
```

| フラグ | 必須 | 説明 |
|---|---|---|
| `--text` | ✅ | 投稿するテキスト |

```
=== Posting to 2 channels ===
  ✓ @myaccount  → post_xxxxxxxx
  ✗ My Page     → failed: channel error
```

#### JSON 出力形式

`--json` フラグを付けると、すべてのコマンドが以下の統一フォーマットで出力します。

```json
{
  "success": true,
  "data": { ... }
}
```

エラー時:

```json
{
  "success": false,
  "error": "エラーメッセージ"
}
```

## 開発

```bash
# ビルド
go build ./...

# テスト（外部 API への通信なし）
go test ./...

# カバレッジ確認
go test -cover ./...
```

## ライセンス

MIT
