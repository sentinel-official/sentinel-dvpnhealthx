package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sentinel-official/sentinel-go-sdk/process"
	"github.com/sentinel-official/sentinelhub/v12/types/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/sentinel-official/sentinel-dvpnhealthx/api"
	"github.com/sentinel-official/sentinel-dvpnhealthx/config"
	"github.com/sentinel-official/sentinel-dvpnhealthx/core"
)

// NewServeCmd returns a new command for serving the healthx API.
func NewServeCmd(cfg *config.Config) *cobra.Command {
	// Default maximum price for a node.
	maxPriceStr := "udvpn:0,13000000"

	// Define the serve command.
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the healthx API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			// Create a new process manager and router.
			manager := process.NewManager(ctx, "manager")
			router := gin.Default()

			// Function to set up everything before starting the server.
			setupFunc := func() error {
				return manager.Setup(func(ctx context.Context) error {
					// Read the home directory from the configuration.
					homeDir := viper.GetString("home")

					// Parse the max price from the string.
					maxPrice, err := v1.NewPriceFromString(maxPriceStr)
					if err != nil {
						return fmt.Errorf("parsing max price %q: %w", maxPriceStr, err)
					}

					// Create a new context and set it up.
					c := core.NewContext()
					if err := c.Setup(ctx, cfg); err != nil {
						return fmt.Errorf("setting up context: %w", err)
					}

					// Set home directory and max price for the context.
					c.WithHomeDir(homeDir)
					c.WithMaxPrice(maxPrice)
					c.Seal()

					// Define the CORS middleware to allow cross-origin requests.
					middlewares := []gin.HandlerFunc{
						cors.New(
							cors.Config{
								AllowAllOrigins: true,
								AllowMethods:    []string{http.MethodGet, http.MethodPost},
							},
						),
					}

					// Apply the CORS middleware to the router.
					router.Use(middlewares...)

					// Register the routes related to node health checks.
					api.RegisterRoutes(router, c)

					return nil
				})
			}

			// Function to start the server.
			startFunc := func() error {
				return manager.Start(func(ctx context.Context) error {
					// Create a TCP listener for the configured address.
					l, err := net.Listen("tcp", cfg.API.GetAddr())
					if err != nil {
						return fmt.Errorf("creating TCP listener on %q: %w", cfg.API.GetAddr(), err)
					}

					// Start serving the HTTP requests on the listener in a goroutine.
					manager.Go(func(ctx context.Context) error {
						if err := http.Serve(l, router); err != nil {
							return fmt.Errorf("serving: %w", err)
						}

						return nil
					})

					// Manage graceful shutdown when the context is cancelled.
					manager.Go(func(ctx context.Context) error {
						defer func() {
							_ = l.Close()
						}()

						select {
						case <-ctx.Done():
							return ctx.Err()
						}
					})

					return nil
				})
			}

			// Wait function to handle the waiting logic after the server is started.
			waitFunc := func() error {
				return manager.Wait(nil)
			}

			// Stop function to handle the stopping logic when the server is shutting down.
			stopFunc := func() error {
				return manager.Stop(nil)
			}

			// Perform setup before starting the server.
			if err := setupFunc(); err != nil {
				return fmt.Errorf("setting up: %w", err)
			}

			// Create an error group for concurrent error handling.
			eg, ctx := errgroup.WithContext(ctx)

			// Goroutine to handle the server start and wait logic.
			eg.Go(func() error {
				if err := startFunc(); err != nil {
					return fmt.Errorf("starting: %w", err)
				}

				if err := waitFunc(); err != nil {
					return fmt.Errorf("waiting: %w", err)
				}

				return nil
			})

			// Goroutine to handle the stop logic when the context is done.
			eg.Go(func() error {
				<-ctx.Done()
				if err := stopFunc(); err != nil {
					return fmt.Errorf("stopping: %w", err)
				}

				return nil
			})

			// Wait for all goroutines to finish.
			if err := eg.Wait(); err != nil {
				return err
			}

			return nil
		},
	}

	// Set flags for the command from the configuration.
	cfg.SetForFlags(cmd.Flags())

	// Define the maxPrice flag which can be passed as a command-line argument.
	cmd.Flags().StringVarP(&maxPriceStr, "max-price", "", maxPriceStr, "specify the maximum price for a node")

	return cmd
}
