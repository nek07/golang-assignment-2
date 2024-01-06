package main

import (
	"context"
	"log"
	"net/http"
	"time"

	//go get go.mongodb.org/mongo-driver/mongo

	_ "context"
	"fmt"
	_ "log"

	_ "github.com/joho/godotenv/autoload"

	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	Name     string `bson:"name"`
	Username string `bson:"username"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

const uri = "mongodb+srv://ataytoleuov:abylai220439@abylaidb.3jyfmar.mongodb.net/"

var client *mongo.Client

func main() {
	uri := uri

	// Create client options
	clientOptions := options.Client().ApplyURI(uri)

	var err error
	client, err = mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}

	// Create context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas!")

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/submit", submitHandler)

	port := 8080
	fmt.Printf("Server is running on http://localhost:%d\n", port)
	err1 := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err1 != nil {
		fmt.Println("Error:", err1)
	}
	// Replace the following connection string with your MongoDB Atlas connection string

	// Perform your MongoDB operations here
	/*collection := client.Database("go-assignment-2").Collection("users")
	user1 := User{Name: "abylai", Email: "abylaiooaoog@gmail.com", Username: "nek07", Password: "12345678"}
	insertResult, err := collection.InsertOne(context.TODO(), user1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	// Disconnect from MongoDB
	err = client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}*/

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "form.html")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			log.Println("Error parsing form:", err)
			return
		}

		user := getData(r.FormValue("name"), r.FormValue("email"), r.FormValue("username"), r.FormValue("password"))
		log.Printf("Received form data: %+v\n", user)

		insertData(user)

		fmt.Fprintln(w, "Data successfully submitted.")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func insertData(u User) {

	collection := client.Database("go-assignment-2").Collection("users")

	insertResult, err := collection.InsertOne(context.TODO(), u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func getData(name string, email string, username string, password string) User {
	user := User{
		Name:     name,
		Email:    email,
		Username: username,
		Password: password,
	}
	return user
}
