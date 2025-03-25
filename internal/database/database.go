package database

import (
	"github.com/adjivas/eir/internal/database/mongodb"

	"github.com/adjivas/eir/pkg/factory"
)

type DbConnector interface {
	HasEquipementStatus() (bool, error)
	CreateEquipementStatus() (err error)
	DropEquipementStatus() (err error)
}

func NewDbConnector() (DbConnector, error) {
	return mongodb.NewMongoDbConnector(factory.EirConfig.Configuration.Mongodb)
}
