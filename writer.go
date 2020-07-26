package crawler

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocsWriter interface {
	WriteDocs([]News) (n int, err error)
	CheckDuplicate(link string) (err error)
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
	client     *mongo.Client
	collection *mongo.Collection
}

func (wtr *mongoDBWriter) Init() error {
	id := "crowlnews"
	pw := "zum123!2#"
	addrTemplate := "mongodb+srv://%s:%s@crowlnews-qzm5a.mongodb.net/test?retryWrites=true&w=majority"
	mongoDBAddr := fmt.Sprintf(addrTemplate, id, pw)
	clientOptions := options.Client().ApplyURI(mongoDBAddr)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	wtr.client = client
	wtr.collection = wtr.client.Database("test").Collection("news")

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

func (wtr *mongoDBWriter) CheckDuplicate(link string) (err error) {

	url := bson.D{{"url", link}}
	var res News
	return wtr.collection.FindOne(context.TODO(), url).Decode(&res)
}

func CheckDuplicateURL(collection *mongo.Collection, res *News, filterUrl bson.D) (err error) {
	return collection.FindOne(context.TODO(), filterUrl).Decode(res)
}

func (wtr *mongoDBWriter) WriteDocs(docs []News) (n int, err error) {

	for _, doc := range docs {
		_, err := wtr.collection.InsertOne(context.TODO(), doc)
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
