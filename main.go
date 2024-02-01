package main

import (
	"context"
	_ "context"
	"encoding/json"
	"fmt"
	"log"
	_ "log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	_ "github.com/eminetto/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"ass3/db"
	"html/template" //end damir
	"strconv"
	"strings" //Damir
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
type Laptop struct {
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	Description string `json:"description"`
	Price       int    `json:"price"`
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
	err = db.AddNewField(ctx, client)
	if err != nil {
		fmt.Println("Error during migration:", err)
		return
	}

	fmt.Println("Migration executed successfully.")

	//server
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/error", errorPageHandler) // Damir end and start
	http.HandleFunc("/crud", crudHandler)
	http.HandleFunc("/getUser", handleGetUser)
	http.HandleFunc("/updateUser", handleUpdateUser)
	http.HandleFunc("/deleteUser", handleDeleteUser)
	http.HandleFunc("/getAllUsers", handleGetAllUsers)

	http.HandleFunc("/products", productsPageHandler)

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
		db.InsertData(user)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
func getData(name string, email string, username string, password string) db.User {
	user := db.User{
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
func productsPageHandler(w http.ResponseWriter, r *http.Request) {
	// page := r.URL.Query().Get("page")
	brands := []string{r.URL.Query().Get("brand")}
	sortBy := r.URL.Query().Get("sort")
	minPrice, err := strconv.Atoi(r.URL.Query().Get("min"))
	maxPrice, err := strconv.Atoi(r.URL.Query().Get("max"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	if page <= 0 {
		page = 1
	}

	if r.URL.Query().Get("min") == "" {
		minPrice = 0
	}
	if r.URL.Query().Get("max") == "" {
		maxPrice = 999999999
	}
	filter := bson.D{}
	db1 := client.Database("go-assignment-2")
	collection1 := db1.Collection("products")
	// query to get all users
	cursor, err := collection1.Find(context.Background(), filter)
	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(context.Background())

	// Iterate through the cursor and print each user

	// brands := []string{"Apple"}

	result, err := db.FindProductsWithFilters(brands, minPrice, maxPrice, sortBy, page)
	if err != nil {
		log.Fatal("Error calling FindProductsWithFilters: %v", err)
	}
	fmt.Println(result)

	// Render the HTML template with the fetched data
	tmpl, err := template.ParseFiles("products.html")
	if err != nil {
		fmt.Println("Error parsing HTML template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// Execute the template with the list of ViewData items
	tmpl.Execute(w, result)
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
	defer cancel()
	id := r.FormValue("userId")
	foundUser, err := db.FindUserByID(ctx, client, "go-assignment-2", "users", id)
	if err != nil {
		fmt.Println("user not found")
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
	defer cancel()
	userIDHex := r.FormValue("updateUserId")
	newUsername := r.FormValue("newUsername")
	var err error = db.UpdateUserUsernameByID(ctx, client, "go-assignment-2", "users", userIDHex, newUsername)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(" username was successfully updated")
	respondWithMessage(w, "updated ofigeno")
}
func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	userIDHex := r.FormValue("deleteUserId")
	var err error = db.DeleteUserByID(ctx, client, "go-assignment-2", "users", userIDHex)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("User was successfully deleted")
	respondWithMessage(w, "Udalen ofigeno")
}
func handleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Get user ID from the request parameters
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	foundUsers, err := db.GetAllUsers(ctx, client, "go-assignment-2", "users")
	if err != nil {
		fmt.Println("user not found")
		return
	}
	log.Printf("Get user result: %+v\n", foundUsers)
	if foundUsers != nil {
		respondWithJSON(w, foundUsers)
	} else {
		respondWithMessage(w, "Users not found")
	}
}
func respondWithMessage(w http.ResponseWriter, msg string) {
	// Respond with an error message in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	// Respond with user data in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
