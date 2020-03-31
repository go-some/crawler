package crawler

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type DocsReader interface {
	ReaderDocs() (n int, err error)
}

/* define mongoDBReader ( */
type mongoDBReader struct {
	client *mongo.Client
}

func (reader *mongoDBReader) Init() error {
	id := os.Getenv("DBID")
	pw := os.Getenv("DBPW")
	addrTemplate := os.Getenv("DBADDR")
	mongoDBAddr := fmt.Sprintf(addrTemplate, id, pw)
	clientOptions := options.Client().ApplyURI(mongoDBAddr)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	reader.client = client

	if err != nil {
		return err
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return err
	}

	fmt.Println("Connected to MongoDB!")

	return nil
}

func (reader *mongoDBReader) Destroy() error {
	err := reader.client.Disconnect(context.TODO())

	if err != nil {
		return err
	}

	fmt.Println("Connection to MongoDB closed.")

	return nil
}

func (reader *mongoDBReader) ReadDocs(filter bson.D, limit int64) (news []*News, err error) {
	collection := reader.client.Database("test").Collection("news")

	var docs []*News

	findOptions := options.Find()
	findOptions.SetLimit(limit)

	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var doc News
		err := cur.Decode(&doc)
		if err != nil {
			log.Fatal(err)
		}

		docs = append(docs, &doc)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return docs, nil
}

func NewMongoDBReader() *mongoDBReader {
	return &mongoDBReader{}
}

/* ) */
