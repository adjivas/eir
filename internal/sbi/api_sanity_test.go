package sbi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/adjivas/eir/pkg/factory"
	"github.com/stretchr/testify/assert"

	"github.com/adjivas/eir/internal/logger"
	util_logger "github.com/free5gc/util/logger"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

)

func setupHttpServer(t *testing.T) *gin.Engine {
	router := util_logger.NewGinWithLogrus(logger.GinLog)
	dataRepositoryGroup := router.Group(factory.EirDrResUriPrefix)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eir := NewMockEIR(ctrl)
	factory.EirConfig = &factory.Config{
		Configuration: &factory.Configuration{
			DbConnectorType: "mongodb",
			Mongodb:         &factory.Mongodb{},
			Sbi: &factory.Sbi{
				BindingIPv4: "127.0.0.1",
				Port:        8000,
			},
		},
	}
	eir.EXPECT().
		Config().
		Return(factory.EirConfig).
		AnyTimes()

	s := NewServer(eir, "")
	AddService(dataRepositoryGroup, s.getEquipementStatusRoutes())
	return router
}

func setupMongoDB() (*mongo.Database) {
	uri := "mongodb://localhost:27017"
 
	// Create a new client and connect to the server
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// Send a ping to confirm a successful connection
	var result bson.M
	db := client.Database("free5gc")
	if err := db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return db
}

func createEquipementStatus(db *mongo.Database) {
	// Create Collection
	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"pei", "equipement_status"},
		"properties": bson.M{
			"pei": bson.M{
				"bsonType": "string",
				"pattern": "^(imei-[0-9]{15}|imeisv-[0-9]{16}|mac((-[0-9a-fA-F]{2}){6})(-untrusted)?|eui((-[0-9a-fA-F]{2}){8}))$",
				"description": "Data type representing the PEI of the UE",
			},
			"supi": bson.M{
				"bsonType": "string",
				"pattern": "^(imsi-[0-9]{5,15}|nai-.+|gci-.+|gli-.+)$",
				"description": "Data type representing the SUPI of the subscriber",
			},
			"gpsi": bson.M{
				"bsonType": "string",
				"pattern": "^(msisdn-[0-9]{5,15}|extid-[^@]+@[^@]+)$",
				"description": "Data type representing the GPSI of the subscriber",
			},
			"equipement_status": bson.M{
				"bsonType": "string",
				"enum": []string{"WHITELISTED", "BLACKLISTED", "GREYLISTED"},
				"description": "Indicates the PEI is white, black or grey listed",
			},
		},
	}
	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}
	opts := options.CreateCollection().SetValidator(validator)

	err := db.CreateCollection(context.TODO(), "policyData.ues.eirData", opts)
	if err != nil {
		panic(err)
	}
}

func TestEIR_Root(t *testing.T) {
	server := setupHttpServer(t)
	reqUri := factory.EirDrResUriPrefix + "/"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	t.Run("EIR Root", func(t *testing.T) {
		require.Equal(t, http.StatusNotImplemented, rsp.Code)
		require.Equal(t, "Hello World!", rsp.Body.String())
	})
}

func TestEIR_EquipementStatus_NotFound(t *testing.T) {
	server := setupHttpServer(t)
	reqUri := factory.EirDrResUriPrefix + "/equipement-status?pei=43&supi=43&gpsi=43"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	t.Run("EquipementStatus", func(t *testing.T) {
		require.Equal(t, http.StatusNotFound, rsp.Code)
		require.Equal(t, "{\"title\":\"Not found\",\"status\":404,\"detail\":\"Supi not found\",\"cause\":\"ERROR_EQUIPMENT_UNKNOWN\"}", rsp.Body.String())
	})

	// t.Run("EquipementStatus", func(t *testing.T) {
	// 	require.Equal(t, http.StatusOK, rsp.Code)
	// 	require.Equal(t, "{\"a\":43}", rsp.Body.String())
	// })
}

func TestEIR_EquipementStatus_DataBaseInsert(t *testing.T) {
	db := setupMongoDB()

	createEquipementStatus(db)
	createEquipementStatus(db)

	defer func() {
		if err := db.Collection("policyData.ues.eirData").Drop(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Success to Insert documents
	// With valid pei and equipement_status
	res, err := db.Collection("policyData.ues.eirData").InsertOne(context.TODO(), bson.M{
		"pei": "imei-012345678901234",
		"equipement_status": "WHITELISTED",
	})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	// With valid pei, equipement_status, supi and gpsi
	res, err = db.Collection("policyData.ues.eirData").InsertOne(context.TODO(), bson.M{
		"pei": "imei-012345678901234",
		"supi": "imsi-208930000000001",
		"gpsi": "msisdn-00000",
		"equipement_status": "WHITELISTED",
	})
	assert.Nil(t, err)
	assert.NotNil(t, res)

	// Fail to Insert documents
	// With missing Pei
	res, err = db.Collection("policyData.ues.eirData").InsertOne(context.TODO(), bson.M{
		"equipement_status": "BLACKLISTED",
	})
	assert.NotNil(t, err)
	assert.Nil(t, res)
	// With missing equipement_status
	res, err = db.Collection("policyData.ues.eirData").InsertOne(context.TODO(), bson.M{
		"pei": "imei-012345678901234",
	})
	assert.NotNil(t, err)
	assert.Nil(t, res)
	// With missing Pei and equipement_status
	res, err = db.Collection("policyData.ues.eirData").InsertOne(context.TODO(), bson.M{
	})
	assert.NotNil(t, err)
	assert.Nil(t, res)
	// With not valid Pei
	res, err = db.Collection("policyData.ues.eirData").InsertOne(context.TODO(), bson.M{
		"pei": "ppppppppppppp-012345678901234",
		"equipement_status": "BLACKLISTED",
	})
	// With not valid equipement_status
	res, err = db.Collection("policyData.ues.eirData").InsertOne(context.TODO(), bson.M{
		"pei": "imei-012345678901234",
		"equipement_status": "PINKLISTED",
	})
	assert.NotNil(t, err)
	assert.Nil(t, res)
	// With not valid supi and gpsi
	res, err = db.Collection("policyData.ues.eirData").InsertOne(context.TODO(), bson.M{
		"pei": "imei-012345678901234",
		"supi": "sssssssssss",
		"gpsi": "ggggggggggg",
		"equipement_status": "BLACKLISTED",
	})
	assert.NotNil(t, err)
	assert.Nil(t, res)
}
