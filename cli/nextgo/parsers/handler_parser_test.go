package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHandlerParser(t *testing.T) {
	dir, _ := os.Getwd()

	parser := HandlerParser{}
	signature, err := parser.Parse(filepath.Join(dir, "testdata", "handler.go"))

	if err != nil {
		t.Fatal(err)
		return
	}

	if signature.PackageName != "handler" {
		t.Errorf("Expected %s as package name got %s", "handler", signature.PackageName)
	}

	if len(signature.Functions) != 2 {
		t.Errorf("Expected to get 2 functions in the signature got %d", len(signature.Functions))
	}
}
