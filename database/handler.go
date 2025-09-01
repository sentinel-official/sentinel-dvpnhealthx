package database

import (
	"context"
	"errors"

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

	return result.Decode(target)
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

	defer cursor.Close(ctx)

	return cursor.All(ctx, target)
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

	defer cursor.Close(ctx)

	return cursor.All(ctx, target)
}

// InsertOne inserts a document into the col.
func (db *handler) InsertOne(ctx context.Context, colName string, document interface{}, opts ...interface{}) (*mongo.InsertOneResult, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	insertOpts := extractOptions[options.InsertOneOptions](opts)
	col := db.Collection(colName, colOpts)

	return col.InsertOne(ctx, document, insertOpts)
}

// UpdateOne updates a single document matching filter.
func (db *handler) UpdateOne(ctx context.Context, colName string, filter bson.M, update bson.M, opts ...interface{}) (*mongo.UpdateResult, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	updateOpts := extractOptions[options.UpdateOptions](opts)
	col := db.Collection(colName, colOpts)

	return col.UpdateOne(ctx, filter, update, updateOpts)
}

// UpdateMany updates multiple documents matching filter.
func (db *handler) UpdateMany(ctx context.Context, colName string, filter bson.M, update bson.M, opts ...interface{}) (*mongo.UpdateResult, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	updateOpts := extractOptions[options.UpdateOptions](opts)
	col := db.Collection(colName, colOpts)

	return col.UpdateMany(ctx, filter, update, updateOpts)
}

// ReplaceOne replaces a document matching filter with replacement.
func (db *handler) ReplaceOne(ctx context.Context, colName string, filter bson.M, replacement interface{}, opts ...interface{}) (*mongo.UpdateResult, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	replaceOpts := extractOptions[options.ReplaceOptions](opts)
	col := db.Collection(colName, colOpts)

	return col.ReplaceOne(ctx, filter, replacement, replaceOpts)
}

// DeleteOne deletes a single document matching filter.
func (db *handler) DeleteOne(ctx context.Context, colName string, filter bson.M, opts ...interface{}) (*mongo.DeleteResult, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	deleteOpts := extractOptions[options.DeleteOptions](opts)
	col := db.Collection(colName, colOpts)

	return col.DeleteOne(ctx, filter, deleteOpts)
}

// DeleteMany deletes multiple documents matching filter.
func (db *handler) DeleteMany(ctx context.Context, colName string, filter bson.M, opts ...interface{}) (*mongo.DeleteResult, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	deleteOpts := extractOptions[options.DeleteOptions](opts)
	col := db.Collection(colName, colOpts)

	return col.DeleteMany(ctx, filter, deleteOpts)
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

	return col.Drop(ctx)
}

// BulkWrite performs a bulk write operation.
func (db *handler) BulkWrite(ctx context.Context, colName string, models []mongo.WriteModel, opts ...interface{}) (*mongo.BulkWriteResult, error) {
	colOpts := extractOptions[options.CollectionOptions](opts)
	bulkOpts := extractOptions[options.BulkWriteOptions](opts)
	col := db.Collection(colName, colOpts)

	return col.BulkWrite(ctx, models, bulkOpts)
}
