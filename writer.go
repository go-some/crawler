package crawler

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
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

func (wtr *mongoDBWriter) Init() error {
	id := os.Getenv("DBID")
	pw := os.Getenv("DBPW")
	addrTemplate := os.Getenv("DBADDR")
	mongoDBAddr := fmt.Sprintf(addrTemplate, id, pw)
	clientOptions := options.Client().ApplyURI(mongoDBAddr)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	wtr.client = client

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

func (wtr *mongoDBWriter) Destroy() error {
	err := wtr.client.Disconnect(context.TODO())

	if err != nil {
		return err
	}

	fmt.Println("Connection to MongoDB closed.")

	return nil
}

func (wtr *mongoDBWriter) WriteDocs(docs []News) (n int, err error) {
	collection := wtr.client.Database("test").Collection("news")

	for _, doc := range docs {
		var res News
		filter := bson.D{{"url", doc.Url}}

		err = collection.FindOne(context.TODO(), filter).Decode(&res)
		if err == nil {
			return 0, fmt.Errorf("Already exist (%s)", doc.Url)
		}

		_, err := collection.InsertOne(context.TODO(), doc)
		if err != nil {
			return 0, err
		}
	}

	return len(docs), nil
}

func NewMongoDBWriter() *mongoDBWriter {
	return &mongoDBWriter{}
}

/* ) */
