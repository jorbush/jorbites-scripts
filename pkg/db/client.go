package db

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DBClient interface {
	AssignBadgeToUser(ctx context.Context, userID string, badgeName string) error
	DeleteBadgeFromUser(ctx context.Context, userID string, badgeName string) error
	GetUserBadges(ctx context.Context, userID string) (badges []string, userName string, err error)
	Close(ctx context.Context) error
}

type MongoClient struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoClient(ctx context.Context, uri string, dbName string) (*MongoClient, error) {
	if uri == "" {
		return nil, errors.New("mongodb connection URI cannot be empty")
	}

	if dbName == "" {
		dbName = ParseDatabaseNameFromURI(uri)
		if dbName == "" {
			dbName = "example"
		}
	}

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	db := client.Database(dbName)

	return &MongoClient{
		client: client,
		db:     db,
	}, nil
}

func (m *MongoClient) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *MongoClient) AssignBadgeToUser(ctx context.Context, userID string, badgeName string) error {
	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID '%s': %w", userID, err)
	}

	collection := m.db.Collection("User")

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "badges", Value: badgeName}}}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user badges: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found with ID '%s'", userID)
	}

	return nil
}

func (m *MongoClient) DeleteBadgeFromUser(ctx context.Context, userID string, badgeName string) error {
	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID '%s': %w", userID, err)
	}

	collection := m.db.Collection("User")

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "badges", Value: badgeName}}}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user badges: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found with ID '%s'", userID)
	}

	return nil
}

func (m *MongoClient) GetUserBadges(ctx context.Context, userID string) (badges []string, userName string, err error) {
	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, "", fmt.Errorf("invalid user ID '%s': %w", userID, err)
	}

	collection := m.db.Collection("User")
	filter := bson.D{{Key: "_id", Value: objID}}

	projection := bson.D{
		{Key: "name", Value: 1},
		{Key: "badges", Value: 1},
	}
	opts := options.FindOne().SetProjection(projection)

	var result struct {
		Name   string   `bson:"name"`
		Badges []string `bson:"badges"`
	}

	err = collection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, "", fmt.Errorf("user not found with ID '%s'", userID)
		}
		return nil, "", fmt.Errorf("failed to fetch user: %w", err)
	}

	name := result.Name
	if name == "" {
		name = "Unnamed User"
	}

	userBadges := result.Badges
	if userBadges == nil {
		userBadges = []string{}
	}

	return userBadges, name, nil
}

func ParseDatabaseNameFromURI(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}

	path := strings.TrimPrefix(u.Path, "/")
	if path == "" {
		return ""
	}

	parts := strings.Split(path, "/")
	return parts[0]
}
