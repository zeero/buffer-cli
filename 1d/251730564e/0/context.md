# Session Context

## User Prompts

### Prompt 1

[Request interrupted by user for tool use]

### Prompt 2

Implement the following plan:

# `buffer post` の --text フラグを位置引数に変更

## Context

現在 `buffer post --text "..."` と書く必要があるが、`--text` ラベルは冗長。
`buffer post "..."` と書けるよう、位置引数（cobra の `args`）に変更する。

## 変更ファイル

`internal/cmd/post.go` のみ。

## 変更内容

### 1. `postText` パッケージ変数と `--text` フラグ定義を削除

```go
// 削除
var postText string

// init() 内も削除
postCmd.Flags().StringVar(&postText, "text", "", "...")
```

### 2. `Run` を `Args` バリデーション付きに変更

cobra の `ExactArgs(1)` で引数が1つ必須であることを強制する。

```go
var postCmd = &...

### Prompt 3

a

