package database

import (
	// "go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/adjivas/eir/internal/database/mongodb"
	"github.com/adjivas/eir/pkg/factory"
)

type DbConnector interface {
	// createEquipementStatus(Mongodb *mongo.Database) (err error)
	dropEquipementStatus()
}

func NewDbConnector() DbConnector {
	return mongodb.NewMongoDbConnector(factory.EirConfig.Configuration.Mongodb)
}
