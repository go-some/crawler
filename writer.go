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
type mongoDBWriter struct{}

func (*mongoDBWriter) WriteDocs(docs []News) (n int, err error) {
	id, pw := os.Getenv("DBID"), os.Getenv("DBPW")
	addrTemplate := os.Getenv("DBADDR")
	fmt.Println(id, pw, addrTemplate)
	mongoDBAddr := fmt.Sprintf(addrTemplate, id, pw)
	fmt.Println("db addr", mongoDBAddr)
	clientOptions := options.Client().ApplyURI(mongoDBAddr)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("test").Collection("news")

	for _, doc := range docs {
		insertResult, err := collection.InsertOne(context.TODO(), doc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	}

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

	return len(docs), nil
}

func NewMongoDBWriter() *mongoDBWriter {
	return &mongoDBWriter{}
}

/* ) */
