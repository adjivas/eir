package mongodb

import (
	"net/http"

	"github.com/adjivas/eir/internal/logger"
	"github.com/adjivas/eir/pkg/factory"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/mongoapi"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	EQUIPMENT_UNKNOWN       = "Data not found"
	EQUIPMENT_UNKNOWN_CAUSE = "DATA_NOT_FOUND"
)

type MongoDbConnector struct {
	*factory.Mongodb
}

func NewMongoDbConnector(mongo *factory.Mongodb) MongoDbConnector {
	return MongoDbConnector{
		Mongodb: mongo,
	}
}

func (m MongoDbConnector) GetDataFromDB(
	collName string, filter bson.M) (
	map[string]interface{}, *models.ProblemDetails,
) {
	data, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		return nil, openapi.ProblemDetailsSystemFailure(err.Error())
	}
	if data == nil {
		return nil, &models.ProblemDetails{
			Title:  EQUIPMENT_UNKNOWN,
			Status: http.StatusNotFound,
			Cause:  EQUIPMENT_UNKNOWN_CAUSE,
		}
	}
	return data, nil
}

func (m MongoDbConnector) GetDataFromDBWithArg(collName string, filter bson.M, strength int) (
	map[string]interface{}, *models.ProblemDetails,
) {
	data, err := mongoapi.RestfulAPIGetOne(collName, filter, strength)
	if err != nil {
		return nil, openapi.ProblemDetailsSystemFailure(err.Error())
	}
	if data == nil {
		logger.ConsumerLog.Errorln("filter: ", filter)
		return nil, &models.ProblemDetails{
			Title:  EQUIPMENT_UNKNOWN,
			Status: http.StatusNotFound,
			Cause:  EQUIPMENT_UNKNOWN_CAUSE,
		}
	}

	return data, nil
}
