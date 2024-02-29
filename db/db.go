package db

import (
	"context"
	"fmt"
	"time"

	"regexp"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://localhost:27017/"

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
}

var client *mongo.Client

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
	ID          string `bson:"_id"`
	Brand       string `bson:"brand"`
	Model       string `bson:"model"`
	Description string `bson:"description"`
	Price       int    `bson:"price"`
	Id          string
}

func InsertData(u User) error {
	collection := client.Database("go-assignment-2").Collection("users")

	// Insert user data into the MongoDB collection
	insertResult, err := collection.InsertOne(context.TODO(), u)
	if err != nil {
		log.WithError(err).Error("Error inserting data into MongoDB")
		return err
	}

	log.WithFields(logrus.Fields{
		"action":    "user_creation",
		"result":    insertResult,
		"timestamp": time.Now().Format(time.RFC3339),
	}).Info("User created successfully")

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
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
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
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
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
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
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
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
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
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
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
