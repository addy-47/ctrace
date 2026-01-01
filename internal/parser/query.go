package parser

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

// FindFunctionDefinition locates a function declaration by name in the given tree.
func FindFunctionDefinition(tree *sitter.Tree, source []byte, funcName string) (*sitter.Node, error) {
	// Query to find all function declarations.
	// We capture the declaration node as @func_decl and the name identifier as @func_name.
	queryStr := `
		(function_declaration
			name: (identifier) @func_name
		) @func_decl
	`

	q, err := sitter.NewQuery([]byte(queryStr), golang.GetLanguage())
	if err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}
	defer q.Close()

	cursor := sitter.NewQueryCursor()
	defer cursor.Close()

	cursor.Exec(q, tree.RootNode())

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		// Filter the match manually by checking the function name text.
		// We need to find both @func_name (to check) and @func_decl (to return).
		var declNode *sitter.Node
		var nameNode *sitter.Node

		for _, capture := range match.Captures {
			captureName := q.CaptureNameForId(capture.Index)
			if captureName == "func_decl" {
				declNode = capture.Node
			} else if captureName == "func_name" {
				nameNode = capture.Node
			}
		}

		// If we found both keys (which we should for valid matches)
		if declNode != nil && nameNode != nil {
			// Extract the content of the name node
			actualName := nameNode.Content(source)
			if actualName == funcName {
				return declNode, nil
			}
		}
	}

	return nil, nil // Not found
}
