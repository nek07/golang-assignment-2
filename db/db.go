package db

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"regexp"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func GetUniqueChatIDDocuments() ([]bson.M, error) {
	collection := Client.Database("go-assignment-2").Collection("Chats")

	// Use MongoDB aggregation to get unique chat_id values
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": "$chat_id",
			},
		},
		{
			"$project": bson.M{
				"chat_id": "$_id",
				"_id":     0,
				"sender":  1,
			},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email       string             `bson:"email" json:"email"`
	Username    string             `bson:"username" json:"username"`
	Password    string             `bson:"password" json:"password"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	AccessToken string             `bson:"access_token" json:"access_token"`
	Role        string             `bson:"role" json:"role"`
}
type Laptop struct {
	ID          string `bson:"_id"`
	Brand       string `bson:"brand"`
	Model       string `bson:"model"`
	Description string `bson:"description"`
	Price       int    `bson:"price"`
	Id          string
}
type Comment struct {
	ID       string `bson:"_id,omitempty"`
	LaptopID string `bson:"laptopId"`
	UserName string `bson:"userName"`
	Text     string `bson:"text"`
	Time     string `bson:"time"`
}
type Recom struct {
	ID   string `bson:"_id,omitempty"`
	Text string `bson:"text"`
}
type ChatMessage struct {
	ChatID    string    `bson:"chat_id"`
	Sender    string    `bson:"sender"`
	Message   string    `bson:"message"`
	Timestamp time.Time `bson:"timestamp"`
}

func InsertMessage(cm ChatMessage) error {

	ctx := context.TODO()
	collection := Client.Database("go-assignment-2").Collection("Chats")
	_, err := collection.InsertOne(ctx, cm)
	if err != nil {
		log.Println("Error storing message:", err)
	}
	return err
}
func InsertData(client *mongo.Client, u User) error {

	collection := client.Database("go-assignment-2").Collection("users")

	// Insert user data into the MongoDB collection
	insertResult, err := collection.InsertOne(context.TODO(), u)
	if err != nil {
		log.Println("Error inserting data into MongoDB:", err)
		return err
	}

	log.WithFields(logrus.Fields{
		"action":    "user_creation",
		"result":    insertResult,
		"timestamp": time.Now().Format(time.RFC3339),
	}).Info("User created successfully")

	return nil
}
func FindUserByEmail(client *mongo.Client, email string) (*User, error) {
	collection := client.Database("go-assignment-2").Collection("users")
	var ctx context.Context
	// Convert the hex string to an ObjectId

	// filter to find the document by its email
	filter := bson.M{"email": email}

	// query
	var user User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		log.WithError(err).Error("Error finding user by email")
		return nil, err
	}
	return &user, nil
}
func UpdateUserUsernameByEmail(client *mongo.Client, email string, field string, value string) error {
	collection := client.Database("go-assignment-2").Collection("users")
	var ctx context.Context
	//  hex string to an ObjectId

	// Specify the filter to find the document by its ID
	filter := bson.M{"email": email}

	// Specify the update to change by $set
	update := bson.M{"$set": bson.M{field: value}}

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
func FindUserByID(ctx context.Context, client *mongo.Client, databaseName, collectionName, userIDHex string) (*User, error) {
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
		log.WithError(err).Error("Error finding user by ID")
		return nil, err
	}
	return &user, nil
}
func GetUserEmails(client *mongo.Client) ([]string, error) {
	// Access the "users" collection in the "your_database" database
	collection := client.Database("go-assignment-2").Collection("users")

	// Set a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the filter to retrieve all documents
	filter := bson.D{{}}

	// Perform the find operation
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Create a slice to store the user emails
	var emails []string

	// Iterate through the cursor and append each email to the slice
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		emails = append(emails, user.Email)
	}

	// Check for errors during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}
func FindUserByToken(client *mongo.Client, token string) (*User, error) {
	collection := client.Database("go-assignment-2").Collection("users")
	var ctx context.Context
	// Convert the hex string to an ObjectId
	filter := bson.M{"access_token": token}

	// filter to find the document by its ID
	// query
	var user User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		log.WithError(err).Error("Error finding user by Token")
		return nil, err
	}
	return &user, nil
}
func UpdateUserUsernameByID(ctx context.Context, client *mongo.Client, databaseName, collectionName, userIDHex string, newUsername string) error {
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
	log.WithFields(logrus.Fields{
		"action":      "user_update",
		"userId":      userIDHex,
		"newUsername": newUsername,
		"timestamp":   time.Now().Format(time.RFC3339),
	}).Info("User information updated successfully")
	return nil
}
func DeleteUserByID(ctx context.Context, client *mongo.Client, databaseName, collectionName, userIDHex string) error {
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
	log.WithFields(logrus.Fields{
		"action":    "user_deletion",
		"userId":    userIDHex,
		"timestamp": time.Now().Format(time.RFC3339),
	}).Info("User deleted successfully")

	return nil
}

func GetAllUsers(ctx context.Context, client *mongo.Client, databaseName, collectionName string) ([]User, error) {
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
	log.WithFields(logrus.Fields{
		"action":    "get_all_users",
		"timestamp": time.Now().Format(time.RFC3339),
	}).Info("Retrieved all users successfully")
	return users, nil
}
func AddNewField(ctx context.Context, client *mongo.Client) error {
	collection := client.Database("go-assignment-2").Collection("users")

	// Adding new default field
	update := bson.M{"$set": bson.M{"minAge": "18"}}
	_, err := collection.UpdateMany(ctx, bson.M{}, update)
	if err != nil {
		return err
	}
	log.WithFields(logrus.Fields{
		"action":    "migration_up",
		"timestamp": time.Now().Format(time.RFC3339),
	}).Info("Migration Up completed successfully.")
	return nil
}
func FindProductsWithFilters(brands []string, minPrice int, maxPrice int, sortBy string, page int) ([]Laptop, int64, error) {
	limit := 10
	skip := (page - 1) * limit

	// Set up MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://damir:CNW6CNosCC9VFPoG@cluster0.qazvzjk.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		return nil, 0, fmt.Errorf("Error creating MongoDB client: %v", err)
	}
	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Set up the database and collection
	database := client.Database("go-assignment-2")
	collection := database.Collection("products")

	// Build the filter based on input parameters
	filter := bson.M{}

	// Add proccer filter
	// if len(proccers) > 0 {
	// 	filter["proccer"] = bson.M{"$in": proccers}
	// }

	// Add brand filter
	if brands[0] == "" {

	}
	if len(brands) > 0 {
		brandRegex := fmt.Sprintf("^%s", regexp.QuoteMeta(brands[0]))
		filter["brand"] = bson.M{"$regex": primitive.Regex{Pattern: brandRegex, Options: ""}}
	}

	// Add price range filter
	filter["price"] = bson.M{"$gte": minPrice, "$lte": maxPrice}
	totalDocuments, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("Error counting documents: %v", err)
	}

	// Build the options for sorting
	options := options.Find()
	switch sortBy {
	case "asc":
		options.SetSort(bson.D{{"price", 1}})
	case "desc":
		options.SetSort(bson.D{{"price", -1}})
	}
	options.SetSkip(int64(skip)).SetLimit(int64(limit))

	// Execute the find operation
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, 0, fmt.Errorf("Error executing find operation: %v", err)
	}
	defer cursor.Close(ctx)

	// Decode the results
	var laptops []Laptop
	err = cursor.All(ctx, &laptops)
	if err != nil {
		return nil, 0, fmt.Errorf("Error decoding results: %v", err)
	}
	log.WithFields(logrus.Fields{
		"action":    "filter_products",
		"brands":    brands,
		"minPrice":  minPrice,
		"maxPrice":  maxPrice,
		"sortBy":    sortBy,
		"page":      page,
		"timestamp": time.Now().Format(time.RFC3339),
	}).Info("User filtered products")
	return laptops, totalDocuments, nil
}

func FindProductById(id string) (products Laptop, err error) {
	ctx := context.TODO()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://damir:CNW6CNosCC9VFPoG@cluster0.qazvzjk.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		return Laptop{}, fmt.Errorf("Error creating MongoDB client: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return Laptop{}, fmt.Errorf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)
	collection := client.Database("go-assignment-2").Collection("products")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Laptop{}, fmt.Errorf("Invalid ID: %v", err)
	}
	filter := bson.M{"_id": objectID}
	result := collection.FindOne(ctx, filter)
	if err := result.Decode(&products); err != nil {
		if err == mongo.ErrNoDocuments {
			return Laptop{}, fmt.Errorf("Product not found")
		}
		return Laptop{}, fmt.Errorf("Error decoding product: %v", err)
	}

	return products, nil
}
func DeleteProduct(id string) error {
	ctx := context.TODO()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://damir:CNW6CNosCC9VFPoG@cluster0.qazvzjk.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		return fmt.Errorf("Error creating MongoDB client: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)
	collection := client.Database("go-assignment-2").Collection("products")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Invalid ID: %v", err)
	}
	filter := bson.M{"_id": objectID}
	collection.DeleteOne(ctx, filter)
	return nil
}
func UpdateProductInDB(id string, brand string, model string, description string, price int) error {
	ctx := context.TODO()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://damir:CNW6CNosCC9VFPoG@cluster0.qazvzjk.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		return fmt.Errorf("Error creating MongoDB client: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("go-assignment-2").Collection("products")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Invalid ID: %v", err)
	}

	filter := bson.M{"_id": objectID}

	update := bson.M{
		"$set": bson.M{
			"brand":       brand,
			"model":       model,
			"description": description,
			"price":       price,
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Error updating product: %v", err)
	}

	return nil
}
func AddProduct(brand string, model string, description string, price int) error {
	ctx := context.TODO()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://damir:CNW6CNosCC9VFPoG@cluster0.qazvzjk.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		return fmt.Errorf("Error creating MongoDB client: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("go-assignment-2").Collection("products")

	product := bson.D{
		{"brand", brand},
		{"model", model},
		{"description", description},
		{"price", price},
	}

	_, err = collection.InsertOne(ctx, product)
	if err != nil {
		return fmt.Errorf("Error inserting product into MongoDB: %v", err)
	}

	return nil
}
func AddComment(username, text, laptopID string) error {
	// Установка соединения с MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://damir:CNW6CNosCC9VFPoG@cluster0.qazvzjk.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		return err
	}
	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	// Получение доступа к коллекции комментариев
	collection := client.Database("go-assignment-2").Collection("comments")

	// Создание нового комментария
	newComment := Comment{
		UserName: username,
		Text:     text,
		LaptopID: laptopID,
		Time:     time.Now().Format(time.RFC822),
	}

	// Вставка нового комментария в базу данных
	_, err = collection.InsertOne(ctx, newComment)
	if err != nil {
		return err
	}

	return nil
}
func GetCommentsByLaptop(laptopID string) ([]Comment, error) {
	// Установка соединения с MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://damir:CNW6CNosCC9VFPoG@cluster0.qazvzjk.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	// Получение доступа к коллекции комментариев
	collection := client.Database("go-assignment-2").Collection("comments")

	// Задание фильтра для поиска комментариев по laptopId
	filter := bson.M{"laptopId": laptopID}

	// Выполнение запроса к базе данных
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Декодирование результатов запроса
	var comments []Comment
	for cursor.Next(ctx) {
		var comment Comment
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
func UpdateAccount(id string, username string, email string) error {
	// Инициализация клиента MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("Error creating MongoDB client: %v", err)
	}

	// Подключение к базе данных
	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Получение коллекции "accounts"
	collection := client.Database("go-assignment-2").Collection("users")

	// Подготовка данных для обновления
	filter := bson.M{"email": id}
	update := bson.M{"$set": bson.M{"username": username, "email": email}}

	// Выполнение обновления
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Error updating account: %v", err)
	}

	return nil
}

func InsertRecom(text string) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("go-assignment-2").Collection("recommendations")

	newComment := Recom{
		Text: text,
	}

	_, err = collection.InsertOne(ctx, newComment)
	if err != nil {
		return err
	}

	return nil
}
