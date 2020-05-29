package test

import (
	"github.com/shun-shun123/bus-timer/src/config"
	"testing"
)

func TestTrimParentheses(t *testing.T) {
	str := "(定時運行)"
	trimStr :=config.TrimParentheses(str)
	if trimStr != "定時運行" {
		t.Fatalf("%vが%vとなっています。", str, trimStr)
	}
}

func TestTrimNoParentheses(t *testing.T) {
	str := "定時運行"
	trimStr := config.TrimParentheses(str)
	if trimStr != "定時運行" {
		t.Fatalf("%vが%vとなっています。", str, trimStr)
	}
}