package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/sentinel-official/sentinel-go-sdk/cmd"
	"github.com/sentinel-official/sentinel-go-sdk/libs/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sentinel-official/sentinel-dvpnhealthx/config"
)

// NewRootCmd creates and returns the root command for the CLI.
func NewRootCmd(userDir string) *cobra.Command {
	// Declare variables for CLI flags
	var (
		homeDir   = filepath.Join(userDir, ".sentinel-dvpnhealthx")
		logFormat = "text"
		logLevel  = "info"
	)

	// Initialize default configuration
	cfg := config.DefaultConfig()

	// Create the root command
	rootCmd := &cobra.Command{
		Use: "sentinel-dvpnhealthx",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// Initialize logger with selected format and level
			logger, err := log.NewLogger(cmd.OutOrStdout(), logFormat, logLevel)
			if err != nil {
				return fmt.Errorf("initializing logger: %w", err)
			}

			// Set the global logger instance
			log.SetLogger(logger)

			// Update the configuration
			cfg.Keyring.HomeDir = homeDir
			cfg.Keyring.Input = cmd.InOrStdin()

			log.Info("Validating configuration")
			if err := cfg.Validate(); err != nil {
				return fmt.Errorf("validating config: %w", err)
			}

			return nil
		},
	}

	// Add subcommands
	rootCmd.AddCommand(
		NewServeCmd(cfg),
		cmd.NewVersionCmd(),
	)

	// Add persistent flags
	rootCmd.PersistentFlags().StringVar(&homeDir, "home", homeDir, "home directory for application config and data")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log.format", logFormat, "format of the log output (json or text)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log.level", logLevel, "log level for output (debug, error, info, none, warn)")

	// Bind flags to global viper instance
	_ = viper.BindPFlag("home", rootCmd.PersistentFlags().Lookup("home"))
	_ = viper.BindPFlag("log.format", rootCmd.PersistentFlags().Lookup("log.format"))
	_ = viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log.level"))

	return rootCmd
}
