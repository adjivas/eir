package util

import (
	"encoding/json"

	"github.com/adjivas/eir/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
)

func ToBsonM(data interface{}) bson.M {
	tmp, err := json.Marshal(data)
	if err != nil {
		logger.UtilLog.Error(err)
	}
	putData := bson.M{}
	err = json.Unmarshal(tmp, &putData)
	if err != nil {
		logger.UtilLog.Error(err)
	}
	return putData
}
