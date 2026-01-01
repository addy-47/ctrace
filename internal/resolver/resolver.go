package resolver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ctrace/internal/parser"
	"ctrace/internal/project"
)

// ResolveFunction defines the location of a function given its import path and name.
// Returns the file path and line number (0-indexed).
func ResolveFunction(importPath string, funcName string) (string, int, error) {
	moduleName, err := project.GetModuleName()
	if err != nil {
		return "", 0, fmt.Errorf("failed to get module name: %w", err)
	}

	// Check if it's an internal package
	if !strings.HasPrefix(importPath, moduleName) {
		return "", 0, nil // External dependency, ignore for now
	}

	// Convert import path to local directory path
	// e.g. "ctrace/internal/cli" -> "./internal/cli"
	relPath := strings.TrimPrefix(importPath, moduleName)
	dirPath := "." + relPath

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".go") || strings.HasSuffix(f.Name(), "_test.go") {
			continue
		}

		filePath := filepath.Join(dirPath, f.Name())
		tree, source, err := parser.ParseFile(filePath)
		if err != nil {
			// Skip files that fail to parse, maybe log valid warning in real app
			continue
		}
		defer tree.Close()

		node, err := parser.FindFunctionDefinition(tree, source, funcName)
		if err != nil {
			continue
		}

		if node != nil {
			return filePath, int(node.StartPoint().Row), nil
		}
	}

	return "", 0, fmt.Errorf("function %s not found in %s", funcName, importPath)
}
