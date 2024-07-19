package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client *mongo.Client
}

type Movie struct {
	ID         string `bson:"_id,omitempty"`
	Email      string `bson:"email"`
	Name       string `bson:"name"`
	LastName   string `bson:"last_name"`
	FirstName  string `bson:"first_name"`
	AvatarURL  string `bson:"avatar_url"`
	NickName   string `bson:"nick_name"`
	ProviderID string `bson:"provider_id"`
}

func Connect(uri string) (*Database, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	return &Database{client: client}, nil
}

func (db *Database) Disconnect() {
	if err := db.client.Disconnect(context.TODO()); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	}
}

func (db *Database) GetMovies() ([]*Movie, error) {
	collection := db.client.Database("moviesdb").Collection("movies")
	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var movies []*Movie
	for cur.Next(context.TODO()) {
		var movie Movie
		err := cur.Decode(&movie)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	return movies, nil
}

func (db *Database) GetMovie(id string) (*Movie, error) {
	collection := db.client.Database("moviesdb").Collection("movies")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var movie Movie
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&movie)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func (db *Database) CreateMovie(movie *Movie) error {
	collection := db.client.Database("moviesdb").Collection("movies")
	_, err := collection.InsertOne(context.TODO(), movie)
	if err != nil {
		return err
	}

	return nil
}
