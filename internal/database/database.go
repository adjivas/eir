package database

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/adjivas/eir/internal/logger"
	"github.com/free5gc/openapi/models"

	"github.com/adjivas/eir/internal/database/mongodb"
	"github.com/adjivas/eir/pkg/factory"
)

const (
	DBCONNECTOR_TYPE_MONGODB factory.DbType = "mongodb"
)

type DbConnector interface {
	GetDataFromDB(collName string, filter bson.M) (map[string]interface{}, *models.ProblemDetails)
	GetDataFromDBWithArg(collName string, filter bson.M, strength int) (map[string]interface{}, *models.ProblemDetails)
}

func NewDbConnector(dbName factory.DbType) DbConnector {
	if dbName == DBCONNECTOR_TYPE_MONGODB {
		return mongodb.NewMongoDbConnector(factory.EirConfig.Configuration.Mongodb)
	} else {
		logger.DbLog.Fatalf("Unsupported database type: %s", dbName)
		return nil
	}
}
