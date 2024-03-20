package clients

import (
	"context"
	"github.com/Braly-Ltd/t2v-api-adapter/properties"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

func NewMongoClient(props *properties.MongoProperties, httpClient *http.Client) (*mongo.Client, error) {
	uri := "mongodb://"
	for _, host := range props.Hosts {
		uri += host + ","
	}
	opts := options.Client().
		SetHTTPClient(httpClient).
		SetConnectTimeout(10 * time.Second).
		ApplyURI(uri[:len(uri)-1])
	return mongo.Connect(context.Background(), opts)
}
