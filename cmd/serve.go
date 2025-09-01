package cmd

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sentinel-official/sentinelhub/v12/types/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sentinel-official/sentinel-dvpnhealthx/api"
	"github.com/sentinel-official/sentinel-dvpnhealthx/config"
	"github.com/sentinel-official/sentinel-dvpnhealthx/core"
)

// NewServeCmd returns a new command for serving the healthx API.
func NewServeCmd(cfg *config.Config) *cobra.Command {
	// Default maximum price for a node
	maxPriceStr := "udvpn:0,13000000"

	// Define the serve command
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the healthx API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			homeDir := viper.GetString("home")

			maxPrice, err := v1.NewPriceFromString(maxPriceStr)
			if err != nil {
				return fmt.Errorf("parsing max price %q: %w", maxPriceStr, err)
			}

			ctx := core.NewContext()
			if err := ctx.Setup(cmd.Context(), cfg); err != nil {
				return fmt.Errorf("setting up context: %w", err)
			}

			ctx.WithHomeDir(homeDir)
			ctx.WithMaxPrice(maxPrice)
			ctx.Seal()

			// Define CORS middleware
			middlewares := []gin.HandlerFunc{
				cors.New(
					cors.Config{
						AllowAllOrigins: true,
						AllowMethods:    []string{http.MethodGet, http.MethodPost},
					},
				),
			}

			// Create a new Gin router and register middleware
			r := gin.Default()
			r.Use(middlewares...)

			// Register API routes for node health checks
			api.RegisterRoutes(r, ctx)

			// Get the API address from the configuration
			addr := cfg.API.GetAddr()
			if err := http.ListenAndServe(addr, r); err != nil {
				return fmt.Errorf("failed to listen and serve: %w", err)
			}

			return nil
		},
	}

	// Set flags for the command from the configuration
	cfg.SetForFlags(cmd.Flags())

	// Define the maxPrice flag
	cmd.Flags().StringVarP(&maxPriceStr, "max-price", "", maxPriceStr, "specify the maximum price for a node")

	return cmd
}
