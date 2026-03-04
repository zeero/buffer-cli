package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
)

// テスト時のみ差し替えて os.Exit を抑制する
var exitFunc = os.Exit

type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func printSuccess(humanOutput func(), data interface{}) {
	if jsonOutput {
		out, _ := json.MarshalIndent(JSONResponse{Success: true, Data: data}, "", "  ")
		fmt.Println(string(out))
	} else {
		humanOutput()
	}
}

func printError(msg string) {
	if jsonOutput {
		out, _ := json.MarshalIndent(JSONResponse{Success: false, Error: msg}, "", "  ")
		fmt.Println(string(out))
	} else {
		color.Red("Error: %s", msg)
	}
	exitFunc(1)
}
