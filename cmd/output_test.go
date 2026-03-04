package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintSuccess_JSONOutput(t *testing.T) {
	jsonOutput = true
	t.Cleanup(func() { jsonOutput = false })

	out := captureStdout(t, func() {
		printSuccess(func() {}, map[string]string{"key": "value"})
	})

	assert.Contains(t, out, `"success": true`)
	assert.Contains(t, out, `"key": "value"`)
}

func TestPrintSuccess_HumanOutput(t *testing.T) {
	jsonOutput = false

	called := false
	captureStdout(t, func() {
		printSuccess(func() { called = true }, nil)
	})

	assert.True(t, called)
}

func TestPrintError_ExitsWithCode1(t *testing.T) {
	origExit := exitFunc
	t.Cleanup(func() { exitFunc = origExit })

	var gotCode int
	exitFunc = func(code int) { gotCode = code }

	jsonOutput = false
	printError("test error")

	assert.Equal(t, 1, gotCode)
}

func TestPrintError_JSONOutput(t *testing.T) {
	origExit := exitFunc
	t.Cleanup(func() { exitFunc = origExit })
	exitFunc = func(code int) {}

	jsonOutput = true
	t.Cleanup(func() { jsonOutput = false })

	out := captureStdout(t, func() {
		printError("something failed")
	})

	assert.Contains(t, out, `"success": false`)
	assert.Contains(t, out, `"something failed"`)
}
