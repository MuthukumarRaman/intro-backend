package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConnObject struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	DbName string `json:"db_name" bson:"db_name"`
	UserId string `json:"user_id" bson:"user_id"`
	Pwd    string `json:"pwd"`
}

var ctx = context.Background()
var MongoClient *mongo.Client
var DBError error

var DBConnections = make(map[string]*mongo.Database)

// By default create shared db connection
var SharedDB *mongo.Database

func GetConnection() *mongo.Database {
	// Get the configuration from environment variables
	config := Config{
		Host:   GetenvStr("MONGO_SHAREDDB_HOST"),
		DbName: GetenvStr("MONGO_SHAREDDB_NAME"),
		UserId: GetenvStr("MONGO_SHAREDDB_USER"),
		Pwd:    GetenvStr("MONGO_SHAREDDB_PASSWORD"),
	}

	// Use a unique key for the database connection, based on the config parameters
	key := fmt.Sprintf("%s-%s", config.Host, config.DbName)

	// Lock to avoid race conditions

	// Check if the connection already exists
	if connection, exists := DBConnections[key]; exists {
		return connection
	}

	// Create a new connection and store it
	connection := CreateDBConnection(config.Host, config.DbName, config.UserId, config.Pwd)
	DBConnections[key] = connection

	// Return the newly created connection
	return connection
}

// CreateDBConnection creates a new MongoDB connection
func CreateDBConnection(host string, dbName string, userId string, pwd string) *mongo.Database {
	// Build the MongoDB connection URL
	dbUrl := fmt.Sprintf("mongodb+srv://%s:%s@%s", userId, pwd, host)

	// Connect to the MongoDB client
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(dbUrl),
	)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil
	}

	// Check the connection
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Printf("DB Ping Error: %v", err)
		log.Fatal(err)
		return nil
	}

	// Return the database
	return client.Database(dbName)
}

type Config struct {
	Host   string
	Port   int
	DbName string
	UserId string
	Pwd    string
}

func Ping() bool {
	DBError = MongoClient.Ping(context.TODO(), nil)
	if DBError != nil {
		// fmt.Println(DBError)
		return false
	}
	return true
}

func GetenvStr(key string) string {
	return os.Getenv(key)
}

func GetenvInt(key string) int {
	s := GetenvStr(key)
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

func CreateDb(host string, port int, dbName string, userid string, pwd string, collectionName string) *mongo.Database {
	dbUrl := fmt.Sprintf("mongodb+srv://%s:%s@%s", userid, pwd, host)
	// dbUrl := fmt.Sprintf("mongodb://%s:%s@%s:%d/?retryWrites=true&authSource=admin&w=majority&authMechanism=SCRAM-SHA-256", userid, pwd, host, port)

	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(dbUrl),
		//.SetAuth(credential),
	)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Printf("DB Ping Error")
		log.Fatal(err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Database(dbName).CreateCollection(ctx, collectionName)
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Println("Database created successfully:", dbName)
	return client.Database(dbName)
}
