package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sentinel-official/sentinel-dvpnhealthx/models"
)

const ColNameNodes = "nodes"

// NodeFindOne retrieves a single Node document.
func NodeFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...interface{}) (*models.Node, error) {
	h := &handler{db}

	var item models.Node
	if err := h.FindOne(ctx, ColNameNodes, filter, &item, opts...); err != nil {
		return nil, err
	}

	return &item, nil
}

// NodeFindAll retrieves multiple Node documents.
func NodeFindAll(ctx context.Context, db *mongo.Database, filter bson.M, opts ...interface{}) ([]models.Node, error) {
	h := &handler{db}

	var items []models.Node
	if err := h.FindAll(ctx, ColNameNodes, filter, &items, opts...); err != nil {
		return nil, err
	}

	return items, nil
}

// NodeAggregate performs an aggregation pipeline on the Node collection.
func NodeAggregate(ctx context.Context, db *mongo.Database, pipeline []bson.M, opts ...interface{}) ([]models.Node, error) {
	h := &handler{db}

	var items []models.Node
	if err := h.Aggregate(ctx, ColNameNodes, pipeline, &items, opts...); err != nil {
		return nil, err
	}

	return items, nil
}

// NodeInsertOne inserts a new Node document.
func NodeInsertOne(ctx context.Context, db *mongo.Database, item *models.Node, opts ...interface{}) (*mongo.InsertOneResult, error) {
	h := &handler{db}

	return h.InsertOne(ctx, ColNameNodes, item, opts...)
}

// NodeUpdateOne updates a single Node document.
func NodeUpdateOne(ctx context.Context, db *mongo.Database, filter bson.M, update bson.M, opts ...interface{}) (*mongo.UpdateResult, error) {
	h := &handler{db}

	return h.UpdateOne(ctx, ColNameNodes, filter, update, opts...)
}

// NodeUpdateMany updates multiple Node documents.
func NodeUpdateMany(ctx context.Context, db *mongo.Database, filter bson.M, update bson.M, opts ...interface{}) (*mongo.UpdateResult, error) {
	h := &handler{db}

	return h.UpdateMany(ctx, ColNameNodes, filter, update, opts...)
}

// NodeReplaceOne replaces a Node document.
func NodeReplaceOne(ctx context.Context, db *mongo.Database, filter bson.M, replacement *models.Node, opts ...interface{}) (*mongo.UpdateResult, error) {
	h := &handler{db}

	return h.ReplaceOne(ctx, ColNameNodes, filter, replacement, opts...)
}

// NodeDeleteOne deletes a single Node document.
func NodeDeleteOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...interface{}) (*mongo.DeleteResult, error) {
	h := &handler{db}

	return h.DeleteOne(ctx, ColNameNodes, filter, opts...)
}

// NodeDeleteMany deletes multiple Node documents.
func NodeDeleteMany(ctx context.Context, db *mongo.Database, filter bson.M, opts ...interface{}) (*mongo.DeleteResult, error) {
	h := &handler{db}

	return h.DeleteMany(ctx, ColNameNodes, filter, opts...)
}

// NodeCount counts the number of Node documents.
func NodeCount(ctx context.Context, db *mongo.Database, filter bson.M, opts ...interface{}) (int64, error) {
	h := &handler{db}

	return h.CountDocuments(ctx, ColNameNodes, filter, opts...)
}

// NodeDropCollection drops the "nodes" collection.
func NodeDropCollection(ctx context.Context, db *mongo.Database, opts ...interface{}) error {
	h := &handler{db}

	return h.DropCollection(ctx, ColNameNodes, opts...)
}

// NodeBulkWrite performs a bulk write operation on the Node collection.
func NodeBulkWrite(ctx context.Context, db *mongo.Database, models []mongo.WriteModel, opts ...interface{}) (*mongo.BulkWriteResult, error) {
	h := &handler{db}

	return h.BulkWrite(ctx, ColNameNodes, models, opts...)
}
