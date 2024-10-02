package mongo

import (
	"fmt"
	"github.com/bccfilkom/career-path-service/internal/pkg/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

func NewMongoInstance() (*mongo.Database, error) {
	clientOptions := options.Client()

	// Fetch username and password from environment variables
	username := env.GetString("MONGO_DB_USERNAME", "")
	password := env.GetString("MONGO_DB_PASSWORD", "")
	host := env.GetString("MONGO_DB_HOST", "localhost")
	port := env.GetString("MONGO_DB_PORT", "27017")
	dbName := env.GetString("MONGO_DB_NAME", "admin")

	// Construct the MongoDB URI with username and password
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s", username, password, host, port, dbName, "admin")

	// Apply the URI to the client options
	clientOptions.ApplyURI(url)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Optionally, check the connection
	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	// Return the database reference
	return client.Database(dbName), nil

}
