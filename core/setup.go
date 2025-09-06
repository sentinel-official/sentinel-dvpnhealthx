package core

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/sentinel-official/sentinel-go-sdk/core"
	"github.com/sentinel-official/sentinel-go-sdk/libs/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/sentinel-dvpnhealthx/config"
)

func (c *Context) SetupDatabase(ctx context.Context, cfg *config.Config) error {
	// Create a new BSON registry and customize the type mapping for MongoDB.
	registry := bson.NewRegistry()
	registry.RegisterTypeMapEntry(bson.TypeDateTime, reflect.TypeOf(time.Time{}))
	registry.RegisterTypeMapEntry(bson.TypeEmbeddedDocument, reflect.TypeOf(bson.M{}))

	// Create a new MongoDB client options instance.
	opts := options.Client().
		SetAppName("").
		ApplyURI(cfg.DB.GetURI()).
		SetRegistry(registry).
		SetMaxPoolSize(0)

	// If a username and password are specified in the configuration, set up authentication for the client.
	if cfg.DB.Username != "" && cfg.DB.Password != "" {
		opts = opts.SetAuth(
			options.Credential{
				Username: cfg.DB.Username,
				Password: cfg.DB.Password,
			},
		)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Connect to the MongoDB database using the client options.
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return fmt.Errorf("connecting to mongodb: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("ping mongodb: %w", err)
	}

	db := client.Database(cfg.DB.Name)

	// Attach the MongoDB database to the context.
	c.WithDatabase(db)
	return nil
}

func (c *Context) SetupClient(_ context.Context, cfg *config.Config) error {
	cc, err := core.NewClientFromConfig(cfg.Config)
	if err != nil {
		return fmt.Errorf("creating client from config: %w", err)
	}

	// Seal the client.
	cc.Seal()

	// Assign the initialized client to the context.
	c.WithClient(cc)
	return nil
}

func (c *Context) Setup(ctx context.Context, cfg *config.Config) error {
	c.WithChainID(cfg.RPC.GetChainID())
	c.WithKeyringBackend(cfg.Keyring.GetBackend())
	c.WithRPCAddr(cfg.RPC.GetAddr())

	log.Info("Setting up blockchain client")
	if err := c.SetupClient(ctx, cfg); err != nil {
		return fmt.Errorf("setting up blockchain client: %w", err)
	}

	log.Info("Setting up database")
	if err := c.SetupDatabase(ctx, cfg); err != nil {
		return fmt.Errorf("setting up database: %w", err)
	}

	return nil
}
