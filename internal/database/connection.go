package database

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb+srv://damir:CNW6CNosCC9VFPoG@cluster0.qazvzjk.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

var Client *mongo.Client
var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
}

func ConnectDB() {
	uri := uri
	file, _ := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)

	// Create client options
	clientOptions := options.Client().ApplyURI(uri)

	var err error
	Client, err = mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}
	log.WithFields(logrus.Fields{
		"action":    "server_access",
		"timestamp": time.Now().Format(time.RFC3339),
	}).Info("Client accessed the server")
	// Create context

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Connect to MongoDB
	err = Client.Connect(ctx)
	if err != nil {
		log.WithError(err).Fatal("Error connecting to MongoDB")
	}

	// Check the connection
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.WithError(err).Fatal("Error pinging MongoDB")
	}

	fmt.Println("Connected to MongoDB Atlas!")
}
