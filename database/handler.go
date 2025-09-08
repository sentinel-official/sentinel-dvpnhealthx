package database

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// extractOptions is a generic helper to pick an option of type T from a variadic slice.
func extractOptions[T any](opts []interface{}) *T {
	for _, opt := range opts {
		if typedOpt, ok := opt.(*T); ok {
			return typedOpt
		}
	}

	return nil
}

// handleError converts mongo.ErrNoDocuments into a nil error.
func handleError(err error) error {
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}

		return err
	}

	return nil
}

// handler manages MongoDB operations for a given collection.
type handler struct {
	*mongo.Database
}

// FindOne retrieves a single document matching filter and decodes it into target.
func (db *handler) FindOne(ctx context.Context, colName string, filter bson.M, target interface{}, opts ...interface{}) error {
	colOpts := extractOptions[options.CollectionOptions](opts)
	findOneOpts := extractOptions[options.FindOneOptions](opts)
	col := db.Collection(colName, colOpts)

	result := col.FindOne(ctx, filter, findOneOpts)
	if err := result.Err(); err != nil {
		return handleError(err)
	}

	if err := result.Decode(target); err != nil {
		return fmt.Errorf("decoding result: %w", err)
	}

	return nil
}

// FindAll retrieves multiple documents matching filter and decodes them into target (pointer to a slice).
func (db *handler) FindAll(ctx context.Context, colName string, filter bson.M, target interface{}, opts ...interface{}) error {
	colOpts := extractOptions[options.CollectionOptions](opts)
	findOpts := extractOptions[options.FindOptions](opts)
	col := db.Collection(colName, colOpts)

	cursor, err := col.Find(ctx, filter, findOpts)
	if err != nil {
		return handleError(err)
	}

	defer func() {
		_ = cursor.Close(ctx)
	}()

	if err := cursor.All(ctx, target); err != nil {
		return fmt.Errorf("iterating results: %w", err)
	}

	return nil
}

// Aggregate runs an aggregation pipeline and decodes the results into target.
func (db *handler) Aggregate(ctx context.Context, colName string, pipeline []bson.M, target interface{}, opts ...interface{}) error {
	colOpts := extractOptions[options.CollectionOptions](opts)
	aggOpts := extractOptions[options.AggregateOptions](opts)
	col := db.Collection(colName, colOpts)

	cursor, err := col.Aggregate(ctx, pipeline, aggOpts)
	if err != nil {
		return handleError(err)
	}

	defer func() {
		_ = cursor.Close(ctx)
	}()

	if err := cursor.All(ctx, target); err != nil {
		return fmt.Errorf("iterating results: %w", err)
	}

	return nil
}

// InsertOne inserts a document into the col.
func (db *handler) InsertOne(ctx context.Context, colName string, document interface{}, opts ...interface{}) (*mongo.InsertOneResult, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	insertOpts := extractOptions[options.InsertOneOptions](opts)
	col := db.Collection(colName, colOpts)

	result, err := col.InsertOne(ctx, document, insertOpts)
	if err != nil {
		return nil, fmt.Errorf("inserting document: %w", err)
	}

	return result, nil
}

// UpdateOne updates a single document matching filter.
func (db *handler) UpdateOne(ctx context.Context, colName string, filter bson.M, update bson.M, opts ...interface{}) (*mongo.UpdateResult, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	updateOpts := extractOptions[options.UpdateOptions](opts)
	col := db.Collection(colName, colOpts)

	result, err := col.UpdateOne(ctx, filter, update, updateOpts)
	if err != nil {
		return nil, fmt.Errorf("updating document: %w", err)
	}

	return result, nil
}

// UpdateMany updates multiple documents matching filter.
func (db *handler) UpdateMany(ctx context.Context, colName string, filter bson.M, update bson.M, opts ...interface{}) (*mongo.UpdateResult, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	updateOpts := extractOptions[options.UpdateOptions](opts)
	col := db.Collection(colName, colOpts)

	result, err := col.UpdateMany(ctx, filter, update, updateOpts)
	if err != nil {
		return nil, fmt.Errorf("updating documents: %w", err)
	}

	return result, nil
}

// ReplaceOne replaces a document matching filter with replacement.
func (db *handler) ReplaceOne(ctx context.Context, colName string, filter bson.M, replacement interface{}, opts ...interface{}) (*mongo.UpdateResult, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	replaceOpts := extractOptions[options.ReplaceOptions](opts)
	col := db.Collection(colName, colOpts)

	result, err := col.ReplaceOne(ctx, filter, replacement, replaceOpts)
	if err != nil {
		return nil, fmt.Errorf("replacing document: %w", err)
	}

	return result, nil
}

// DeleteOne deletes a single document matching filter.
func (db *handler) DeleteOne(ctx context.Context, colName string, filter bson.M, opts ...interface{}) (*mongo.DeleteResult, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	deleteOpts := extractOptions[options.DeleteOptions](opts)
	col := db.Collection(colName, colOpts)

	result, err := col.DeleteOne(ctx, filter, deleteOpts)
	if err != nil {
		return nil, fmt.Errorf("deleting document: %w", err)
	}

	return result, nil
}

// DeleteMany deletes multiple documents matching filter.
func (db *handler) DeleteMany(ctx context.Context, colName string, filter bson.M, opts ...interface{}) (*mongo.DeleteResult, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	deleteOpts := extractOptions[options.DeleteOptions](opts)
	col := db.Collection(colName, colOpts)

	result, err := col.DeleteMany(ctx, filter, deleteOpts)
	if err != nil {
		return nil, fmt.Errorf("deleting documents: %w", err)
	}

	return result, nil
}

// CountDocuments returns the count of documents matching filter.
func (db *handler) CountDocuments(ctx context.Context, colName string, filter bson.M, opts ...interface{}) (int64, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	countOpts := extractOptions[options.CountOptions](opts)
	col := db.Collection(colName, colOpts)

	count, err := col.CountDocuments(ctx, filter, countOpts)
	if err != nil {
		return 0, handleError(err)
	}

	return count, nil
}

// DropCollection drops the entire col.
func (db *handler) DropCollection(ctx context.Context, colName string, opts ...interface{}) error {
	colOpts := extractOptions[options.CollectionOptions](opts)
	col := db.Collection(colName, colOpts)

	if err := col.Drop(ctx); err != nil {
		return fmt.Errorf("dropping collection: %w", err)
	}

	return nil
}

// BulkWrite performs a bulk write operation.
func (db *handler) BulkWrite(ctx context.Context, colName string, models []mongo.WriteModel, opts ...interface{}) (*mongo.BulkWriteResult, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	bulkOpts := extractOptions[options.BulkWriteOptions](opts)
	col := db.Collection(colName, colOpts)

	result, err := col.BulkWrite(ctx, models, bulkOpts)
	if err != nil {
		return nil, fmt.Errorf("bulk writing: %w", err)
	}

	return result, nil
}
