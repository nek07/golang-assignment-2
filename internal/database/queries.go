package database

import (
	"context"
	"fmt"
	"regexp"
	"store/internal/models"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func InsertMessage(cm models.ChatMessage) error {

	collection := Client.Database("go-assignment-2").Collection("Chats")
	_, err := collection.InsertOne(context.TODO(), cm)
	if err != nil {
		log.Println("Error storing message:", err)
	}
	return err
}
func InsertData(u models.User) error {
	collection := Client.Database("go-assignment-2").Collection("users")data
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
func FindUserByEmail(email string) (*models.User, error) {
	collection := Client.Database("go-assignment-2").Collection("users")
	var ctx context.Context
	// Convert the hex string to an ObjectId

	// filter to find the document by its email
	filter := bson.M{"email": email}

	// query
	var user models.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		log.WithError(err).Error("Error finding user by email")
		return nil, err
	}
	return &user, nil
}
func UpdateUserUsernameByEmail(email string, field string, value string) error {
	collection := Client.Database("go-assignment-2").Collection("users")
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
func GetUserEmails() ([]string, error) {
	collection := Client.Database("go-assignment-2").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var emails []string

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		emails = append(emails, user.Email)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}

func FindUserByToken(token string) (*models.User, error) {
	ctx := context.TODO()
	collection := Client.Database("go-assignment-2").Collection("users")

	filter := bson.M{"access_token": token}
	var user models.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		log.WithError(err).Error("Error finding user by Token")
		return nil, err
	}
	return &user, nil
}
func UpdateUserUsernameByID(userIDHex string, newUsername string) error {
	ctx := context.TODO()
	collection := Client.Database("go-assignment-2").Collection("users")

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
func DeleteUserByID(userIDHex string) error {
	ctx := context.TODO()
	collection := Client.Database("go-assignment-2").Collection("users")

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

func GetAllUsers() ([]models.User, error) {
	ctx := context.TODO()
	collection := Client.Database("go-assignment-2").Collection("users")

	filter := bson.D{}

	// query to get all users
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and print each user
	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
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

func FindProductsWithFilters(brands []string, minPrice int, maxPrice int, sortBy string, page int) ([]models.Laptop, int64, error) {
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
	var laptops []models.Laptop
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

func FindProductById(id string) (products models.Laptop, err error) {
	ctx := context.TODO()

	collection := Client.Database("go-assignment-2").Collection("products")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Laptop{}, fmt.Errorf("Invalid ID: %v", err)
	}
	filter := bson.M{"_id": objectID}
	result := collection.FindOne(ctx, filter)
	if err := result.Decode(&products); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Laptop{}, fmt.Errorf("Product not found")
		}
		return models.Laptop{}, fmt.Errorf("Error decoding product: %v", err)
	}

	return products, nil
}
func DeleteProduct(id string) error {
	ctx := context.TODO()

	collection := Client.Database("go-assignment-2").Collection("products")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %v", err)
	}
	filter := bson.M{"_id": objectID}
	collection.DeleteOne(ctx, filter)
	return nil
}
func UpdateProductInDB(id string, brand string, model string, description string, price int) error {
	ctx := context.TODO()

	collection := Client.Database("go-assignment-2").Collection("products")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %v", err)
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
		return fmt.Errorf("error updating product: %v", err)
	}

	return nil
}
func AddProduct(brand string, model string, description string, price int) error {
	// Create a context for the operation
	ctx := context.TODO()

	// Access the "products" collection in the MongoDB database
	collection := Client.Database("go-assignment-2").Collection("products")

	// Create a BSON document representing the product
	product := bson.D{
		{"brand", brand},
		{"model", model},
		{"description", description},
		{"price", price},
	}

	// Insert the product document into the collection
	_, err := collection.InsertOne(ctx, product)
	if err != nil {
		// If an error occurs during insertion, return an error
		return fmt.Errorf("error inserting product into MongoDB: %v", err)
	}

	// Return nil (no error) if the insertion was successful
	return nil
}

func AddComment(username, text, laptopID string) error {
	ctx := context.TODO()

	// Получение доступа к коллекции комментариев
	collection := Client.Database("go-assignment-2").Collection("comments")

	// Создание нового комментария
	newComment := models.Comment{
		UserName: username,
		Text:     text,
		LaptopID: laptopID,
		Time:     time.Now().Format(time.RFC822),
	}

	// Вставка нового комментария в базу данных
	_, err := collection.InsertOne(ctx, newComment)
	if err != nil {
		// If an error occurs during insertion, return an error
		return fmt.Errorf("error inserting product into MongoDB: %v", err)
	}

	return nil
}
func GetCommentsByLaptop(laptopID string) ([]models.Comment, error) {
	ctx := context.TODO()

	// Получение доступа к коллекции комментариев
	collection := Client.Database("go-assignment-2").Collection("comments")

	// Задание фильтра для поиска комментариев по laptopId
	filter := bson.M{"laptopId": laptopID}

	// Выполнение запроса к базе данных
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Декодирование результатов запроса
	var comments []models.Comment
	for cursor.Next(ctx) {
		var comment models.Comment
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
	// Create a context for the operation
	ctx := context.TODO()

	// Access the "users" collection in the MongoDB database
	collection := Client.Database("go-assignment-2").Collection("users")

	// Prepare the filter to find the user by email
	filter := bson.M{"email": id}

	// Prepare the update operation to set the new username and email
	update := bson.M{"$set": bson.M{"username": username, "email": email}}

	// Execute the update operation
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		// If an error occurs during the update, return an error
		return fmt.Errorf("error updating account: %v", err)
	}

	// Return nil (no error) if the update was successful
	return nil
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
