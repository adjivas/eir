package database

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/free5gc/openapi/models"
	"github.com/adjivas/eir/internal/database/mongodb"
	"github.com/adjivas/eir/internal/logger"
	"github.com/adjivas/eir/pkg/factory"
)

const (
	APPDATA_INFLUDATA_DB_COLLECTION_NAME       = "applicationData.influenceData"
	APPDATA_INFLUDATA_SUBSC_DB_COLLECTION_NAME = "applicationData.influenceData.subsToNotify"
	APPDATA_PFD_DB_COLLECTION_NAME             = "applicationData.pfds"

	DBCONNECTOR_TYPE_MONGODB factory.DbType = "mongodb"
)

type DbConnector interface {
	PatchDataToDBAndNotify(collName string, ueId string, patchItem []models.PatchItem, filter bson.M) (
		map[string]interface{}, map[string]interface{}, error)
	GetDataFromDB(collName string, filter bson.M) (map[string]interface{}, *models.ProblemDetails)
	GetDataFromDBWithArg(collName string, filter bson.M, strength int) (map[string]interface{}, *models.ProblemDetails)
	DeleteDataFromDB(collName string, filter bson.M)
}

func NewDbConnector(dbName factory.DbType) DbConnector {
	if dbName == DBCONNECTOR_TYPE_MONGODB {
		return mongodb.NewMongoDbConnector(factory.EirConfig.Configuration.Mongodb)
	} else {
		logger.DbLog.Fatalf("Unsupported database type: %s", dbName)
		return nil
	}
}
