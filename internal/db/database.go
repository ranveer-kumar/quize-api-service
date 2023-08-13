package db

import (
	"context"
	"fmt"

	// "log"

	log "quize-api-service/internal/logger"
	"quize-api-service/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
)

func InitMongoDB(uri, dbName string) {
	clientOptions := options.Client().ApplyURI(uri)

	// Set up connection pooling and timeout options
	clientOptions.SetMaxPoolSize(50) // Adjust the max pool size as needed
	clientOptions.SetMinPoolSize(5)  // Adjust the min pool size as needed
	clientOptions.SetMaxConnIdleTime(30 * time.Second)
	clientOptions.SetConnectTimeout(5 * time.Second)
	clientOptions.SetSocketTimeout(1 * time.Minute)

	var err error
	mongoClient, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		// log.Fatalf("Error connecting to MongoDB: %v", err)
		log.Error("Error connecting to MongoDB", err)
	}

	// Check the connection
	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		// log.Fatalf("Error pinging MongoDB: %v", err)
		log.Error("Error pinging MongoDB", err)
	}
	// log.Println("Connected to MongoDB")
	log.Info("Connected to MongoDB")

	// mongoClient = client
}

func GetMongoClient() *mongo.Client {
	if mongoClient == nil {
		log.Error("MongoDB client not initialized", nil)
	}
	return mongoClient
}

func GetCollection(collectionName string) *mongo.Collection {
	return GetMongoClient().Database("quize-db").Collection(collectionName)
}

func CloseMongoDB() {
	if mongoClient != nil {
		err := mongoClient.Disconnect(context.Background())
		if err != nil {
			log.Error("Error disconnecting from MongoDB: ", err)
		} else {
			log.Info("Disconnected from MongoDB")
		}
	}
}

func InsertQuize(quize *models.Quize) (string, error) {
	collection := GetCollection("quizes")

	// Create a context with a timeout for the insert operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertResult, err := collection.InsertOne(ctx, quize)
	if err != nil {
		// log.Printf("Error inserting quize: %v", err)
		log.Error("Error inserting quize: ", err)
		return "", err
	}

	insertedID := insertResult.InsertedID.(string)
	return insertedID, err
}

func UpdateQuize(id string, updatedQuize *models.Quize) error {
	collection := GetCollection("quizes")

	// Create a context with a timeout for the query operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{{"$set", updatedQuize}},
	)
	return err
}

func GetQuizeByID(id string) (*models.Quize, error) {
	collection := GetCollection("quizes") // Get the collection instance

	// Create a context with a timeout for the query operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Define a filter to find the quiz by ID
	filter := bson.M{"_id": id}

	var quize models.Quize
	err := collection.FindOne(ctx, filter).Decode(&quize)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Quiz not found
		}
		log.Error("Error fetching quiz by ID: ", err)
		return nil, fmt.Errorf("error fetching quiz by ID")
	}
	return &quize, nil // Found quiz
}

func ListQuizes(pageSize int, pageNumber int) ([]models.Quize, error) {
	collection := GetCollection("quizes")

	// Calculate skip value for pagination
	skip := (pageNumber - 1) * pageSize

	// Create a context with a timeout for the query operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Define options for the query (e.g., sorting and pagination)
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"createdAt", -1}}) // Sort by createdAt field in descending order
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))

	// Execute the query
	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Error("Error querying quizes collection: ", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var quizes []models.Quize
	for cursor.Next(ctx) {
		var quize models.Quize
		if err := cursor.Decode(&quize); err != nil {
			log.Error("Error decoding quize: ", err)
			return nil, err
		}
		// Convert the time fields to IST
		quize.ConvertTimeToIST()
		quizes = append(quizes, quize)
	}

	return quizes, nil
}

func GetTotalQuizeCount() (int, error) {
	collection := GetCollection("quizes")

	// Create a context with a timeout for the query operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Execute the count query
	totalCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Error("Error counting quizzes: ", err)
		return 0, err
	}

	return int(totalCount), nil
}

func DeleteQuize(id string) error {
	collection := GetCollection("quizes")

	// Create a context with a timeout for the delete operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Define a filter to find the quiz by ID
	filter := bson.M{"_id": id}

	// Execute the delete operation
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Error("Error deleting quiz: ", err)
		return err
	}

	return nil
}
