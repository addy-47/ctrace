package parser

import (
	"fmt"
	"path"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

// ImportMap maps a package alias (or name) to its full import path.
type ImportMap map[string]string

// ParseImports extracts import declarations from the file's root node.
func ParseImports(rootNode *sitter.Node, source []byte) (ImportMap, error) {
	queryStr := `
		(import_spec
			name: (package_identifier)? @alias
			path: (interpreted_string_literal) @path
		)
	`
	q, err := sitter.NewQuery([]byte(queryStr), golang.GetLanguage())
	if err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}
	defer q.Close()

	cursor := sitter.NewQueryCursor()
	defer cursor.Close()

	cursor.Exec(q, rootNode)

	imports := make(ImportMap)

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		var alias, importPath string

		for _, capture := range match.Captures {
			captureName := q.CaptureNameForId(capture.Index)
			content := capture.Node.Content(source)

			if captureName == "alias" {
				alias = content
			} else if captureName == "path" {
				// Remove quotes from interpreted string literal
				importPath = strings.Trim(content, "\"")
			}
		}

		// If alias is empty, derive it from the last segment of the path
		if alias == "" && importPath != "" {
			alias = path.Base(importPath)
		}

		if alias != "" && importPath != "" {
			imports[alias] = importPath
		}
	}

	return imports, nil
}
