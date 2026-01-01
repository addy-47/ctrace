package cli

import (
	"os"

	"ctrace/internal/config"
	"ctrace/internal/utils"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ctrace",
	Short: "Static Code Lifecycle Tracer",
	Long: `ctrace is a high-performance, deterministic static code analysis CLI 
that traces the complete lifecycle of an endpoint or function across a large codebase.`,
	// Requirement: Disable the default completion command
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Only print banner if no arguments are provided
		utils.PrintBanner()

		// It's good practice to show help if no args provided for a root command that has subcommands.
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig, utils.InitLogger)
}
