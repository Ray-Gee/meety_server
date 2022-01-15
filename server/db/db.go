package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	// "os"

	// "log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/mongodb/mongo-go-driver/bson/primitive"
	// "github.com/mongodb/mongo-go-driver/mongo"
)

type Person struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string `json:"firstname,omitempty" bson:"firstname,omitempty"`
}

// type demo struct {
// 	ID 		primitive.ObjectID		`bson:"_id, omitempty"`
// 	AuthorID 	string				`bson:"author_id"`
// 	Content 	string				`bson:"content"`
// 	Title 		string				`bson:"title"`
// }

var client *mongo.Client

func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	// response.Header().Add("content-type", "application/json")
	// log.Println("request: ", request)
	response.Header().Set("content-type", "application/json")
	var person Person
	b, _ := ioutil.ReadAll(request.Body)
	log.Println("person: ", string(b))

	
	e := json.NewDecoder(request.Body).Decode(&person)
	if e != nil {
		log.Println(e)
		return
	}
	collection := client.Database("meety").Collection("person")
	log.Println("request.Body: ", request.Body)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	log.Println("before: ", person)
	result, err := collection.InsertOne(ctx, person)
	// result, err := collection.InsertOne(ctx, Person{Firstname: string(b)})
	// result, err := collection.InsertOne(ctx, bson.D{
	// 	{Key: "Firstname", Value: "RYUICHI"},
	// })
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message0": "` + err.Error() + `"}`))
		return
	}
	// log.Println("result: ", result)
	json.NewEncoder(response).Encode(result)
}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var people []Person
	collection := client.Database("meety").Collection("person")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message1": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message2": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

func Get(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("ここは動きます!\n"))
}

func main() {

	//Connect to MongoDB
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://ryuichi:C7YiYbiQsHRQruJ@cluster0.odnhq.mongodb.net/meety?retryWrites=true&w=majority"))
	// if err != nil { log.Fatal(err) }

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb+srv://ryuichi:C7YiYbiQsHRQruJ@cluster0.odnhq.mongodb.net/meety?retryWrites=true&w=majority")
	client, _ = mongo.Connect(ctx, clientOptions)
	// if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
	// 	// Can't connect to Mongo server
	// 	log.Fatal(err)
	// }
	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// err = client.Connect(ctx)
	// log.Println("client: ", client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// if err != nil { log.Fatal(err) }
	// client, _ = mongo.Connect(ctx, "mongodb://localhost:27017")

	// collection := client.Database("meety").Collection("person")
	// log.Println("collection: ", *collection)
	// srv := &http.Server{
	// 	Addr:    fmt.Sprintf("%s:%d", os.Getenv("HOSTNAME"), os.Getenv("PORT")),
	// }
	router := mux.NewRouter()
	port := ":12345"
	host := "localhost"
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/get", Get).Methods("GET")
    fmt.Printf("DB server started at: %s\n", host + port)
	http.ListenAndServe(host + port, router)

	// Close connection
	client.Disconnect(ctx)
}