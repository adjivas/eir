package sbi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	db "github.com/adjivas/eir/internal/database"
	"github.com/adjivas/eir/internal/logger"
	"github.com/adjivas/eir/pkg/factory"
	util_logger "github.com/free5gc/util/logger"
	"github.com/free5gc/util/mongoapi"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

)

type testdata struct {
	influId string
	supi    string
}

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
	dataRepositoryRoutes := s.getDataRepositoryRoutes()
	AddService(dataRepositoryGroup, dataRepositoryRoutes)
	return router
}

func setupMongoDB(t *testing.T) {
	err := mongoapi.SetMongoDB("test5gc", "mongodb://localhost:27017")
	require.Nil(t, err)
	err = mongoapi.Drop(db.APPDATA_INFLUDATA_DB_COLLECTION_NAME)
	require.Nil(t, err)
	err = mongoapi.Drop(db.APPDATA_INFLUDATA_SUBSC_DB_COLLECTION_NAME)
	require.Nil(t, err)
	err = mongoapi.Drop(db.APPDATA_PFD_DB_COLLECTION_NAME)
	require.Nil(t, err)
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

func TestEIR_EquipementStatus_Adjivas(t *testing.T) {
	uri := "mongodb://localhost:27017"
	// Use the SetServerAPIOptions() method to set the Stable API version to 1

	// Create a new client and connect to the server
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	var result bson.M
	db := client.Database("free5gc")
	if err := db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	// Create Collection
	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"name", "age"},
		"properties": bson.M{
			"name": bson.M{
				"bsonType": "string",
				"enum": []string{"alex", "la"},
				"description": "the name of the user, which is required and " +
					"must be a string",
			},
			"age": bson.M{
				"bsonType": "int",
				"minimum":  18,
				"description": "the age of the user, which is required and " +
					"must be an integer >= 18",
			},
		},
	}
	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}
	opts := options.CreateCollection().SetValidator(validator)

	err = db.CreateCollection(context.TODO(), "policyData.ues.eirData", opts)
	if err != nil {
		panic(err)
	}

	// Insert One Document
	res, err := db.Collection("policyData.ues.eirData").InsertOne(context.TODO(), bson.M{"name": "alex", "age": 31})
	if err != nil {
		panic(err)
	}
	fmt.Printf("inserted document with ID %v\n", res.InsertedID)

	// Drop Collection
	err = db.Collection("policyData.ues.eirData").Drop(context.TODO())
	if err != nil {
		panic(err)
	}
}
