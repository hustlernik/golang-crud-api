package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// Declare global variables for the MongoDB client and collection
var client *mongo.Client
var collection *mongo.Collection

// Initialize MongoDB connection
func initMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://hustlernik99:il3qobkV3bonxhbb@cluster0.txatf.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	// Access the collection
	collection = client.Database("mydb").Collection("users")
	fmt.Println("Connected to MongoDB!")
	return nil
}

type User struct {
  Name     string `json:"name"`
	Email    string `json:"email"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "application/json")
		var user User
		err:= json.NewDecoder(r.Body).Decode(&user)
		if err!= nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }
		result, err := collection.InsertOne(context.Background(), user)
		if err!= nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
		fmt.Println(result)

}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	var user User
	err:= collection.FindOne(context.Background(),bson.M{"name":name}).Decode(&user)
	if err!=nil {
		fmt.Println("can not find name")
	}

	w.Header().Set("content-type","application/json")

   if err := json.NewEncoder(w).Encode(user.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}



}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{name}",getUser).Methods("GET")

	if err := initMongoDB(); err != nil {
		fmt.Printf("Error initializing MongoDB: %v\n", err)
		return
	}

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			fmt.Printf("Error disconnecting from MongoDB: %v\n", err)
		}
	}()

  fmt.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
