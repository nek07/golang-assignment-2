package main

import (
	// "context"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	//go get go.mongodb.org/mongo-driver/mongo

	_ "context"
	"fmt"
	_ "log"

	_ "github.com/joho/godotenv/autoload"

	_ "github.com/eminetto/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	Username   string             `bson:"username"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	Created_at time.Time
	Updated_at time.Time
}

const uri = "mongodb://localhost:27017/"

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

	// Create context

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
	err = addNewField(ctx, client)
	if err != nil {
		fmt.Println("Error during migration:", err)
		return
	}

	fmt.Println("Migration executed successfully.")
	//get and show all users
	// users, err := getAllUsers(ctx, client, "go-assignment-2", "users")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println("All Users:")
	// for _, user := range users {
	// 	fmt.Printf("ID: %s, Username: %s\n", user.ID, user.Username)
	// }

	//find user by ID
	// userIDHex := "65a96fd970c38547a91d4db3"
	// user, err := findUserByID(ctx, client, "go-assignment-2", "users", userIDHex)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(user)

	//update username by ID
	// newUsername := "Dosyan"
	// err = updateUserUsernameByID(ctx, client, "go-assignment-2", "users", userIDHex, newUsername)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(" username was successfully updated")

	// //delete document by id (commented because, we need to write always new id, cause of err)
	// userIDHexDeletion := "65a79472c03f5ac2c93660fd"
	// err = deleteUserByID(ctx, client, "go-assignment-2", "users", userIDHexDeletion)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	//server
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/error", errorPageHandler) // Damir end and start
	http.HandleFunc("/crud", crudHandler)
	http.HandleFunc("/getUser", handleGetUser)
	http.HandleFunc("/updateUser", handleUpdateUser)
	// http.HandleFunc("/deleteUser", handleDeleteUser)
	// http.HandleFunc("/getAllUsers", handleGetAllUsers)

	port := 8080
	fmt.Printf("Server is running on http://localhost:%d\n", port)
	err1 := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err1 != nil {
		fmt.Println("Error:", err1)
	}

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
		insertData(user)
		fmt.Fprintln(w, "Data successfully submitted.")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func insertData(u User) {
	collection := client.Database("go-assignment-2").Collection("users")

	// Insert user data into the MongoDB collection
	insertResult, err := collection.InsertOne(context.TODO(), u)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}
func getData(name string, email string, username string, password string) User {
	user := User{
		ID:         primitive.NewObjectID(),
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
func crudHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "crud.html")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func handleGetUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from the request parameters
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel();
	id := r.FormValue("userId");
	foundUser, err := findUserByID(ctx, client, "go-assignment-2", "users", id)
	if err != nil{
		fmt.Println("user not found");
		return
	}
	log.Printf("Get user result: %+v\n", foundUser)
    // Convert userID to int
    

    // Find user by ID (dummy data for illustration)
    
    // Respond with user data in a JSON format
    if foundUser != nil {
        respondWithJSON(w, foundUser)
    } else {
        respondWithMessage(w, "User not found")
    }
}
func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel();
	userIDHex := r.FormValue("updateUserId");
	newUsername := r.FormValue("newUsername");
	var err error = updateUserUsernameByID(ctx, client, "go-assignment-2", "users", userIDHex, newUsername)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(" username was successfully updated")
	respondWithMessage(w, "updated ofigeno")
}

func respondWithMessage(w http.ResponseWriter, msg string) {
    // Respond with an error message in JSON format
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNotFound)
    json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
    // Respond with user data in JSON format
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
func findUserByID(ctx context.Context, client *mongo.Client, databaseName, collectionName, userIDHex string) (*User, error) {
	collection := client.Database(databaseName).Collection(collectionName)

	// Convert the hex string to an ObjectId
	objectID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return nil, err
	}

	// filter to find the document by its ID
	filter := bson.M{"_id": objectID}

	// query
	var user User
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
func updateUserUsernameByID(ctx context.Context, client *mongo.Client, databaseName, collectionName, userIDHex string, newUsername string) error {
	collection := client.Database(databaseName).Collection(collectionName)

	//  hex string to an ObjectId
	objectID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return err
	}

	// Specify the filter to find the document by its ID
	filter := bson.M{"_id": objectID}

	// Specify the update to change by $set
	update := bson.M{"$set": bson.M{"username": newUsername}}

	// Update query
	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if updateResult.ModifiedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
func deleteUserByID(ctx context.Context, client *mongo.Client, databaseName, collectionName, userIDHex string) error {
	collection := client.Database(databaseName).Collection(collectionName)

	// Convert the hex string to an ObjectId
	objectID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return err
	}

	// filter to find the document by its ID
	filter := bson.M{"_id": objectID}

	// deletion
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if deleteResult.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func getAllUsers(ctx context.Context, client *mongo.Client, databaseName, collectionName string) ([]User, error) {
	collection := client.Database(databaseName).Collection(collectionName)

	filter := bson.D{}

	// query to get all users
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and print each user
	var users []User
	for cursor.Next(ctx) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)

		// Print user details
		fmt.Printf("ID: %s, Username: %s\n", user.ID, user.Username)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
func addNewField(ctx context.Context, client *mongo.Client) error {
	collection := client.Database("go-assignment-2").Collection("users")

	// Adding new default field
	update := bson.M{"$set": bson.M{"minAge": "18"}}
	_, err := collection.UpdateMany(ctx, bson.M{}, update)
	if err != nil {
		return err
	}

	fmt.Println("Migration Up completed successfully.")
	return nil
}
