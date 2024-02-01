package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://localhost:27017/"

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
	Brand       string `bson:"brand"`
	Model       string `bson:"model"`
	Description string `bson:"description"`
	Price       int    `bson:"price"`
}

func InsertData(u User) {
	collection := client.Database("go-assignment-2").Collection("users")

	// Insert user data into the MongoDB collection
	insertResult, err := collection.InsertOne(context.TODO(), u)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

/*
	func GetProducts(ctx context.Context, client *mongo.Client,databaseName, collectionName) ([]laptops,error) {
		cursor, err := collection1.Find(ctx, filter)
		if err != nil {
			fmt.Println(err)
		}
		defer cursor.Close(ctx)

		// Iterate through the cursor and print each user
		var laptops []Laptop
		for cursor.Next(ctx) {
			var laptop Laptop
			err := cursor.Decode(&laptop)
			if err != nil {
				fmt.Println(err)
			}
			laptops = append(laptops, laptop)

		}
		fmt.Println(len(laptops))
		if err := cursor.Err(); err != nil {
			fmt.Println(err)
		}
	}
*/
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

	fmt.Println("Migration Up completed successfully.")
	return nil
}
func FindProductsWithFilters(brands []string, minPrice int, maxPrice int, sortBy string, page int) ([]Laptop, int64, error) {
	limit := 10
	skip := (page - 1) * limit

	// Set up MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil,0,  fmt.Errorf("Error creating MongoDB client: %v", err)
	}
	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		return nil,0,  fmt.Errorf("Error connecting to MongoDB: %v", err)
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
		return nil,0,  fmt.Errorf("Error executing find operation: %v", err)
	}
	defer cursor.Close(ctx)

	// Decode the results
	var laptops []Laptop
	err = cursor.All(ctx, &laptops)
	if err != nil {
		return nil,0, fmt.Errorf("Error decoding results: %v", err)
	}

	return laptops, totalDocuments, nil
}
