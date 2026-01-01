package parser

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

// ExtractCalls analyzes a function node and returns a list of called function names.
func ExtractCalls(funcNode *sitter.Node, source []byte) ([]string, error) {
	queryStr := `
		(call_expression
			function: (_) @callee
		)
	`

	q, err := sitter.NewQuery([]byte(queryStr), golang.GetLanguage())
	if err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}
	defer q.Close()

	cursor := sitter.NewQueryCursor()
	defer cursor.Close()

	// Execute the query only on the provided function node
	cursor.Exec(q, funcNode)

	var calls []string
	seen := make(map[string]bool)

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			captureName := q.CaptureNameForId(capture.Index)
			if captureName == "callee" {
				// The captured node is the function expression (identifier or selector)
				name := capture.Node.Content(source)
				if !seen[name] {
					calls = append(calls, name)
					seen[name] = true
				}
			}
		}
	}

	return calls, nil
}
