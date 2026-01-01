package cli

import (
	"fmt"
	"strings"

	"ctrace/internal/parser"
	"ctrace/internal/resolver"
	"ctrace/internal/utils"

	"github.com/spf13/cobra"
)

// explainCmd represents the explain command
var explainCmd = &cobra.Command{
	Use:   "explain [file-path] [function-name]",
	Short: "Explain the flow of a given endpoint or function",
	Long: `Traces and explains the lifecycle of the provided entry point.
For Phase 2B testing, this command finds a specific function definition in the file.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		funcName := args[1]

		utils.Log.Info(fmt.Sprintf("Parsing file: %s...", filePath))

		tree, source, err := parser.ParseFile(filePath)
		if err != nil {
			utils.Log.Error(fmt.Sprintf("Error parsing file: %v", err))
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer tree.Close()

		// extract imports
		imports, err := parser.ParseImports(tree.RootNode(), source)
		if err != nil {
			utils.Log.Error(fmt.Sprintf("Error parsing imports: %v", err))
			// Don't error out, just log
		} else {
			fmt.Println("Imports found:")
			if len(imports) == 0 {
				fmt.Println("  - (none)")
			} else {
				for alias, path := range imports {
					fmt.Printf("  - %s -> %s\n", alias, path)
				}
			}
		}

		utils.Log.Info(fmt.Sprintf("Searching for function: %s...", funcName))
		node, err := parser.FindFunctionDefinition(tree, source, funcName)
		if err != nil {
			utils.Log.Error(fmt.Sprintf("Error searching function: %v", err))
			fmt.Printf("Error: %v\n", err)
			return
		}

		if node != nil {
			// Tree-sitter rows are 0-indexed.
			fmt.Printf("Found function '%s' at line %d to %d\n", funcName, node.StartPoint().Row, node.EndPoint().Row)

			utils.Log.Info("Extracting calls...")
			calls, err := parser.ExtractCalls(node, source)
			if err != nil {
				utils.Log.Error(fmt.Sprintf("Error extracting calls: %v", err))
				fmt.Printf("Error extracting calls: %v\n", err)
				return
			}

			fmt.Printf("Calls made by %s:\n", funcName)
			if len(calls) == 0 {
				fmt.Println("  - (none)")
			} else {
				for _, call := range calls {
					// Check for resolver logic
					// e.g. "cli.Execute" -> alias="cli", fn="Execute"
					display := call
					if strings.Contains(call, ".") {
						parts := strings.Split(call, ".")
						if len(parts) == 2 {
							alias := parts[0]
							fn := parts[1]

							if pkgPath, ok := imports[alias]; ok {
								resolvedFile, line, err := resolver.ResolveFunction(pkgPath, fn)
								if err == nil && resolvedFile != "" {
									// Resolution successful
									// Using ANSI Green code for the resolved part
									display = fmt.Sprintf("%s \033[32m--> defined in %s:%d\033[0m", call, resolvedFile, line+1) // +1 for human readable 1-indexing
								}
							}
						}
					}
					fmt.Printf("  - %s\n", display)
				}
			}
		} else {
			fmt.Printf("Function '%s' not found in %s\n", funcName, filePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
