package crawler

import (
	"context"
	"fmt"
	"log"
	"os"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocsWriter interface {
	WriteDocs([]News) (n int, err error)
}

/* define printerWriter ( */
type printerWriter struct{}

func (*printerWriter) WriteDocs(docs []News) (n int, err error) {
	fmt.Println(docs)
	return len(docs), nil
}

func NewPrinterWriter() *printerWriter {
	return &printerWriter{}
}

/* ) */

/* define mongodbWriter ( */
type mongoDBWriter struct {
	client *mongo.Client
}

func (wtr *mongoDBWriter) Init() {
	// [TODO] err -> panic?
	id := os.Getenv("DBID")
	pw := os.Getenv("DBPW")
	addrTemplate := os.Getenv("DBADDR")
	mongoDBAddr := fmt.Sprintf(addrTemplate, id, pw)
	clientOptions := options.Client().ApplyURI(mongoDBAddr)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	wtr.client = client

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}

func (wtr *mongoDBWriter) Destroy() {
	// [TODO] err -> panic?
	err := wtr.client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to MongoDB closed.")
}

func (wtr *mongoDBWriter) WriteDocs(docs []News) (n int, err error) {
	collection := wtr.client.Database("test").Collection("news")

	for _, doc := range docs {
		insertResult, err := collection.InsertOne(context.TODO(), doc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	}

	return len(docs), nil
}

func NewMongoDBWriter() *mongoDBWriter {
	return &mongoDBWriter{}
}

/* ) */
