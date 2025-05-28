package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConfig holds MongoDB connection parameters from environment variables.
type MongoConfig struct {
	Username     string `env:"MONGODB_USERNAME"`
	Password     string `env:"MONGODB_PASSWORD"`
	Host         string `env:"MONGODB_HOST,required"`
	Port         int    `env:"MONGODB_PORT" envDefault:"27017"`
	DatabaseName string `env:"MONGODB_DATABASE_NAME,required"`
	IsSRV        bool   `env:"MONGODB_IS_SRV" envDefault:"false"`
	AuthSource   string `env:"MONGODB_AUTH_SOURCE" envDefault:"admin"` // Common default, or can be DatabaseName
}

// ConnectDB establishes a connection to MongoDB.
// It returns the client, a disconnect function, and the parsed config.
func ConnectDB(cfg *MongoConfig) *mongo.Client {
	var mongoURI string
	scheme := "mongodb"
	if cfg.IsSRV {
		scheme = "mongodb+srv"
	}

	hostPort := cfg.Host
	if !cfg.IsSRV && cfg.Port != 0 { // Port is not used for SRV records
		hostPort = cfg.Host + ":" + strconv.Itoa(cfg.Port)
	}

	userInfo := ""
	if cfg.Username != "" {
		userInfo = url.UserPassword(cfg.Username, cfg.Password).String() + "@"
	}

	// Basic URI structure - don't include database name in the path when using authentication
	// The database name in the path is optional and mainly used for the default database
	mongoURI = fmt.Sprintf("%s://%s%s/", scheme, userInfo, hostPort)

	// Add query parameters like authSource
	queryParams := url.Values{}
	if cfg.Username != "" && cfg.AuthSource != "" { // Only add authSource if username is present
		queryParams.Add("authSource", cfg.AuthSource)
	}
	// Add other necessary SRV parameters if any (e.g. replicaSet, but often handled by SRV record itself)
	if cfg.IsSRV {
		queryParams.Add("retryWrites", "true")
		queryParams.Add("w", "majority")
	}

	if len(queryParams) > 0 {
		mongoURI += "?" + queryParams.Encode()
	}

	log.Printf("Constructed MongoDB URI: %s", strings.Replace(mongoURI, cfg.Password, "****", 1))

	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the primary
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v (URI used: %s)", err, strings.Replace(mongoURI, cfg.Password, "****", 1))
	}

	fmt.Printf("Successfully connected to MongoDB! (Database: %s)\n", cfg.DatabaseName)

	return client
}
