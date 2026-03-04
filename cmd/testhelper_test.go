package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// newMockGraphQLServer は data を {"data": ...} でラップして返す httptest サーバーを作成する
func newMockGraphQLServer(t *testing.T, data interface{}) *httptest.Server {
	t.Helper()
	body, err := json.Marshal(map[string]interface{}{"data": data})
	if err != nil {
		t.Fatalf("failed to marshal mock data: %v", err)
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}))
}

// newMockGraphQLErrorServer は GraphQL errors フィールドを含むレスポンスを返す httptest サーバーを作成する
func newMockGraphQLErrorServer(t *testing.T, msg string) *httptest.Server {
	t.Helper()
	body, _ := json.Marshal(map[string]interface{}{
		"errors": []map[string]string{{"message": msg}},
	})
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}))
}

// newMockServer は正常 / エラーのどちらかの httptest サーバーを作成する
func newMockServer(t *testing.T, data interface{}, errMsg string) *httptest.Server {
	t.Helper()
	if errMsg != "" {
		return newMockGraphQLErrorServer(t, errMsg)
	}
	return newMockGraphQLServer(t, data)
}

// setupTest はテスト用グローバル変数をセットし、終了時に元の値へ戻す
func setupTest(t *testing.T, endpoint string) {
	t.Helper()
	origEndpoint := bufferEndpoint
	origToken := apiToken
	origOrg := orgFlag
	origJSON := jsonOutput

	bufferEndpoint = endpoint
	apiToken = "test-token"

	t.Cleanup(func() {
		bufferEndpoint = origEndpoint
		apiToken = origToken
		orgFlag = origOrg
		jsonOutput = origJSON
	})
}

// captureStdout は os.Stdout をパイプに差し替えて fn 内の出力文字列を返す
// fatih/color の colored 出力（color.Output）はキャプチャ対象外
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = orig
	out, _ := io.ReadAll(r)
	return string(out)
}
