package main

import (
	// "context"
	"context"
	"log"
	"net/http"
	"time"

	//go get go.mongodb.org/mongo-driver/mongo

	_ "context"
	"fmt"
	_ "log"

	_ "github.com/joho/godotenv/autoload"

	_ "github.com/eminetto/mongo-migrate"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	// "go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"

	"html/template" //end damir
	"strings"       //Damir
)

type User struct {
	Name       string `bson:"name"`
	Username   string `bson:"username"`
	Email      string `bson:"email"`
	Password   string `bson:"password"`
	Created_at time.Time
	Updated_at time.Time
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	// Define the migration
	// Perform your migration task, e.g., add an index or update documents

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/error", errorPageHandler) // Damir end and start

	port := 8080
	fmt.Printf("Server is running on http://localhost:%d\n", port)
	err1 := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err1 != nil {
		fmt.Println("Error:", err1)
	}
	// Replace the following connection string with your MongoDB Atlas connection string

	// Perform your MongoDB operations here

	// Disconnect from MongoDB
	err = client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { //Damir
		error404PageHandler(w, r)
		return
	} // end damir

	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "form.html")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/submit" { //Damir
		error404PageHandler(w, r)
		return
	} // end damir

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			log.Println("Error parsing form:", err)
			return
		}
		user := getData(r.FormValue("name"), r.FormValue("email"), r.FormValue("username"), r.FormValue("password"))

		errors := checkForm(user.Name, user.Email, user.Username, user.Password, r.FormValue("confirm-password")) //Damir
		if errors.NameError != "" || errors.EmailError != "" || errors.UsernameError != "" ||
			errors.PasswordError != "" || errors.ConfirmPasswordError != "" {
			tmpl, err := template.ParseFiles("error.html")
			if err != nil {
				http.Error(w, "Error rendering error page", http.StatusInternalServerError)
				log.Println("Error rendering error page:", err)
				return
			}

			err = tmpl.Execute(w, errors)
			if err != nil {
				http.Error(w, "Error rendering error page", http.StatusInternalServerError)
				log.Println("Error rendering error page:", err)
				return
			}
			return
		} // end damir

		log.Printf("Received form data: %+v\n", user)

		// insertData(user)

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
		Name:       name,
		Email:      email,
		Username:   username,
		Password:   password,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}
	return user
}

// Damir
type ValidationErrors struct {
	NameError            string
	EmailError           string
	UsernameError        string
	PasswordError        string
	ConfirmPasswordError string
}

func checkForm(name, email, username, password, confirmPassword string) ValidationErrors {
	var errors ValidationErrors

	// Name validation
	if strings.TrimSpace(name) == "" {
		errors.NameError = "Name is required."
	}

	// Email validation
	if strings.TrimSpace(email) == "" {
		errors.EmailError = "Email is required."
	} else if !strings.Contains(email, "@") {
		errors.EmailError = "Invalid email address."
	}

	// Username validation
	if strings.TrimSpace(username) == "" {
		errors.UsernameError = "Username is required."
	}

	// Password validation
	if strings.TrimSpace(password) == "" {
		errors.PasswordError = "Password is required."
	} else if len(password) < 8 {
		errors.PasswordError = "Password must be at least 8 characters long."
	}

	// Confirm Password validation
	if password != confirmPassword {
		errors.ConfirmPasswordError = "Passwords do not match."
	}

	return errors
}
func errorPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "error.html")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func error404PageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "error404.html")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
} //end damir
