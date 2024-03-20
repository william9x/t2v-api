package properties

import (
	"github.com/golibs-starter/golib/config"
)

type MongoProperties struct {
	Hosts []string // Format: [user:pwd@host1:port1, user:pwd@host2:port2, ...]
}

func NewMongoProperties(loader config.Loader) (*MongoProperties, error) {
	props := MongoProperties{}
	err := loader.Bind(&props)
	return &props, err
}

func (r *MongoProperties) Prefix() string {
	return "app.mongo"
}
