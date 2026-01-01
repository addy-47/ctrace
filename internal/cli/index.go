package cli

import (
	"ctrace/internal/utils"
	"fmt"

	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index the codebase for analysis",
	Long:  `Scans and indexes the codebase to build the necessary symbol maps and call graphs.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Log.Info("Indexing codebase...")
		fmt.Println("Indexing codebase...")
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
