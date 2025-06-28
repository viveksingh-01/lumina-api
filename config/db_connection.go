package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/viveksingh-01/lumina-api/handlers"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectToDB() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal("Error occurred while connecting to the database: ", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Could not establish connection to the database.")
	}
	log.Println("Connected to database successfully!")

	DB := client.Database("luminadb")
	if DB == nil {
		log.Fatal("Database connection is not initialized")
	}
}
