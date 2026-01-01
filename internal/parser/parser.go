package parser

import (
	"context"
	"fmt"
	"io/ioutil"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

// GetParser returns a configured tree-sitter parser for the specified language.
// Currently only supports "go".
func GetParser(lang string) (*sitter.Parser, error) {
	if lang != "go" {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}

	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())
	return parser, nil
}

// ParseFile reads a file and returns its parsed AST tree and the raw content.
// The caller is responsible for Close()ing the tree if necessary, though
// in this CLI context it will be reclaimed on exit.
func ParseFile(filePath string) (*sitter.Tree, []byte, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	parser, err := GetParser("go")
	if err != nil {
		return nil, nil, err
	}
	defer parser.Close()

	tree, err := parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse file content: %w", err)
	}

	return tree, content, nil
}
