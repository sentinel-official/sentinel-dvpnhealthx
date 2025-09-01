package node

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sentinel-official/sentinel-go-sdk/libs/log"
	"github.com/sentinel-official/sentinel-go-sdk/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"

	"github.com/sentinel-official/sentinel-dvpnhealthx/core"
	"github.com/sentinel-official/sentinel-dvpnhealthx/database"
	"github.com/sentinel-official/sentinel-dvpnhealthx/models"
)

func handlerGetNode(c *core.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req, err := NewRequestGetNode(ctx)
		if err != nil {
			err := fmt.Errorf("parsing request from context: %w", err)
			ctx.JSON(http.StatusBadRequest, types.NewResponseError(1, err))
			return
		}

		filter := bson.M{
			"addr": req.URI.Addr,
		}
		projection := bson.M{
			"_id":         0,
			"location.ip": 0,
		}
		opts := options.FindOne().
			SetProjection(projection)

		item, err := database.NodeFindOne(ctx, c.Database(), filter, opts)
		if err != nil {
			err := fmt.Errorf("querying node from database: %w", err)
			ctx.JSON(http.StatusInternalServerError, types.NewResponseError(2, err))
			return
		}

		ctx.JSON(http.StatusOK, types.NewResponseResult(item))
	}
}

func handlerGetNodes(c *core.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := NewRequestGetNodes(ctx)
		if err != nil {
			err := fmt.Errorf("parsing request from context: %w", err)
			ctx.JSON(http.StatusBadRequest, types.NewResponseError(1, err))
			return
		}

		filter := bson.M{}
		projection := bson.M{
			"_id":         0,
			"location.ip": 0,
		}
		sort := bson.D{
			{Key: "timestamp", Value: -1},
		}
		opts := options.Find().
			SetProjection(projection).
			SetSort(sort)

		items, err := database.NodeFindAll(ctx, c.Database(), filter, opts)
		if err != nil {
			err := fmt.Errorf("querying nodes from database: %w", err)
			ctx.JSON(http.StatusInternalServerError, types.NewResponseError(2, err))
			return
		}

		ctx.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}

func handlerInspectNode(c *core.Context) gin.HandlerFunc {
	setup := func(ctx context.Context, addr string) error {
		filter := bson.M{
			"addr": addr,
		}

		item, err := database.NodeFindOne(ctx, c.Database(), filter)
		if err != nil {
			return fmt.Errorf("querying node from database: %w", err)
		}
		if !item.IsZero() {
			return nil
		}

		keyAddr, err := c.UpsertKey(ctx, addr)
		if err != nil {
			return fmt.Errorf("upserting key %q: %w", addr, err)
		}
		if err := c.UpsertGrants(ctx, keyAddr); err != nil {
			return fmt.Errorf("upserting grants for addr %q: %w", keyAddr.String(), err)
		}

		item = &models.Node{Addr: addr}
		if _, err := database.NodeInsertOne(ctx, c.Database(), item); err != nil {
			return fmt.Errorf("inserting node in database: %w", err)
		}

		return nil
	}

	m := &sync.Map{}
	return func(ctx *gin.Context) {
		req, err := NewRequestInspectNode(ctx)
		if err != nil {
			err := fmt.Errorf("parsing request from context: %w", err)
			ctx.JSON(http.StatusBadRequest, types.NewResponseError(1, err))
			return
		}

		_, loaded := m.LoadOrStore(req.Body.Addr, &struct{}{})
		if loaded {
			err := fmt.Errorf("inspection for node %q is already in progress", req.Body.Addr)
			ctx.JSON(http.StatusConflict, types.NewResponseError(2, err))
			return
		}

		go func(ctx context.Context, addr string) {
			defer m.Delete(addr)

			filter := bson.M{
				"addr": addr,
			}
			item := &models.Node{
				Addr:      addr,
				Timestamp: time.Now().Unix(),
			}

			eg := &errgroup.Group{}
			eg.Go(func() error {
				log.Info("Setting up inspection", "addr", addr)
				if err := setup(ctx, addr); err != nil {
					return fmt.Errorf("setting up inspection: %w", err)
				}

				log.Info("Inspecting node", "addr", addr)

				output, err := c.Inspect(ctx, addr)
				if err != nil {
					return fmt.Errorf("inspecting node: output %q: %w", string(output), err)
				}

				if err := json.Unmarshal(output, &item.Location); err != nil {
					return fmt.Errorf("unmarshaling inpection output: %w", err)
				}

				return nil
			})

			if err := eg.Wait(); err != nil {
				item.Error = err.Error()
			}

			log.Info("Inspection completed", "addr", addr, "error", item.Error)

			if _, err := database.NodeReplaceOne(ctx, c.Database(), filter, item); err != nil {
				log.Error("Replacing node in database", "addr", addr, "cause", err)
				return
			}
		}(context.TODO(), req.Body.Addr)

		ctx.JSON(http.StatusCreated, types.NewResponseResult(nil))
		return
	}
}
