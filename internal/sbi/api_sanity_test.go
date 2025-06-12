package sbi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adjivas/eir/internal/logger"
	"github.com/adjivas/eir/internal/sbi/processor"
	"github.com/adjivas/eir/internal/util"
	"github.com/adjivas/eir/pkg/factory"
	eir_api_service "github.com/free5gc/openapi/eir/EIRService"
	"github.com/free5gc/openapi/models"
	util_logger "github.com/free5gc/util/logger"
	"github.com/free5gc/util/mongoapi"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func setupHttpServer(t *testing.T) *gin.Engine {
	return setupHttpServerWithDefaultStatus(t, "")
}

func setupHttpServerWithDefaultStatus(t *testing.T, defaultStatus string) *gin.Engine {
	router := util_logger.NewGinWithLogrus(logger.GinLog)
	equipmentStatusGroup := router.Group(factory.EirDrResUriPrefix)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eir := NewMockEIR(ctrl)
	configuration := factory.Configuration{
		DbConnectorType: "mongodb",
		Mongodb:         &factory.Mongodb{},
		Sbi: &factory.Sbi{
			BindingIP: "127.0.0.1",
			Port:      8000,
		},
		DefaultStatus: defaultStatus,
	}
	if defaultStatus != "" {
		configuration.DefaultStatus = defaultStatus
	}
	factory.EirConfig = &factory.Config{
		Configuration: &configuration,
	}
	eir.EXPECT().
		Config().
		Return(factory.EirConfig).
		AnyTimes()

	processor := processor.NewProcessor(eir)
	eir.EXPECT().Processor().Return(processor).AnyTimes()

	s := NewServer(eir, "")
	equipmentStatusRoutes := s.getEquipmentStatusRoutes()
	AddService(equipmentStatusGroup, equipmentStatusRoutes)
	return router
}

func setupMongoDB(t *testing.T) {
	err := mongoapi.SetMongoDB("test5gc", "mongodb://localhost:27017")
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

func TestEIR_EquipmentStatus_FoundEquipmentStatus(t *testing.T) {
	server := setupHttpServer(t)
	setupMongoDB(t)

	defer func() {
		if err := mongoapi.Drop("policyData.ues.eirData"); err != nil {
			panic(err)
		}
	}()

	filter := bson.M{"pei": nil}
	pei1 := bson.M{"pei": "imei-42", "equipment_status": "BLACKLISTED"}
	pei2 := bson.M{"pei": "imei-43", "equipment_status": "BLACKLISTED"}
	pei3 := bson.M{"pei": "imei-012345678901234", "equipment_status": "WHITELISTED"}
	filters := []bson.M{filter, filter, filter}
	peis := []map[string]interface{}{pei1, pei2, pei3}
	err := mongoapi.RestfulAPIPutMany("policyData.ues.eirData", filters, peis)
	assert.Nil(t, err)

	reqUri := factory.EirDrResUriPrefix + "/equipment-status?pei=imei-012345678901234"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	expected_message := util.ToBsonM(eir_api_service.EIREquipmentStatusGetResponse{
		Status: "WHITELISTED",
	})
	t.Run("EquipmentStatus", func(t *testing.T) {
		json_message := eir_api_service.EIREquipmentStatusGetResponse{}

		err := json.Unmarshal(rsp.Body.Bytes(), &json_message)
		assert.Nil(t, err)

		message := util.ToBsonM(json_message)

		require.Equal(t, expected_message, message)
		require.Equal(t, http.StatusOK, rsp.Code)
	})
}

func TestEIR_EquipmentStatus_FoundEquipmentStatus_WithSUPI(t *testing.T) {
	server := setupHttpServer(t)
	setupMongoDB(t)

	defer func() {
		if err := mongoapi.Drop("policyData.ues.eirData"); err != nil {
			panic(err)
		}
	}()

	filter := bson.M{"pei": nil}
	pei1 := bson.M{"pei": "imei-012345678901234", "supi": "imsi-208930000000001", "equipment_status": "BLACKLISTED"}
	pei2 := bson.M{"pei": "imei-43", "supi": "imsi-208930123456789", "equipment_status": "BLACKLISTED"}
	pei3 := bson.M{"pei": "imei-012345678901234", "supi": "imsi-208930123456789", "equipment_status": "WHITELISTED"}
	filters := []bson.M{filter, filter, filter}
	peis := []map[string]interface{}{pei1, pei2, pei3}
	err := mongoapi.RestfulAPIPutMany("policyData.ues.eirData", filters, peis)
	assert.Nil(t, err)

	reqUri := factory.EirDrResUriPrefix + "/equipment-status?pei=imei-012345678901234&supi=imsi-208930123456789"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	expected_message := util.ToBsonM(eir_api_service.EIREquipmentStatusGetResponse{
		Status: "WHITELISTED",
	})
	t.Run("EquipmentStatus", func(t *testing.T) {
		json_message := eir_api_service.EIREquipmentStatusGetResponse{}

		err := json.Unmarshal(rsp.Body.Bytes(), &json_message)
		assert.Nil(t, err)

		message := util.ToBsonM(json_message)

		require.Equal(t, expected_message, message)
		require.Equal(t, http.StatusOK, rsp.Code)
	})
}

func TestEIR_EquipmentStatus_FoundEquipmentStatus_WithGPSI(t *testing.T) {
	server := setupHttpServer(t)
	setupMongoDB(t)

	defer func() {
		if err := mongoapi.Drop("policyData.ues.eirData"); err != nil {
			panic(err)
		}
	}()

	filter := bson.M{"pei": nil}
	pei1 := bson.M{"pei": "imei-42", "gpsi": "msisdn-00042", "equipment_status": "BLACKLISTED"}
	pei2 := bson.M{"pei": "imei-012345678901234", "gpsi": "msisdn-00000", "equipment_status": "BLACKLISTED"}
	pei3 := bson.M{"pei": "imei-012345678901234", "gpsi": "msisdn-12345", "equipment_status": "WHITELISTED"}
	filters := []bson.M{filter, filter, filter}
	peis := []map[string]interface{}{pei1, pei2, pei3}
	err := mongoapi.RestfulAPIPutMany("policyData.ues.eirData", filters, peis)
	assert.Nil(t, err)

	reqUri := factory.EirDrResUriPrefix + "/equipment-status?pei=imei-012345678901234&gpsi=msisdn-12345"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	expected_message := util.ToBsonM(eir_api_service.EIREquipmentStatusGetResponse{
		Status: "WHITELISTED",
	})
	t.Run("EquipmentStatus", func(t *testing.T) {
		json_message := eir_api_service.EIREquipmentStatusGetResponse{}

		err := json.Unmarshal(rsp.Body.Bytes(), &json_message)
		assert.Nil(t, err)

		message := util.ToBsonM(json_message)

		require.Equal(t, expected_message, message)
		require.Equal(t, http.StatusOK, rsp.Code)
	})
}

func TestEIR_EquipmentStatus_FoundEquipmentStatus_WithSUPI_GPSI(t *testing.T) {
	server := setupHttpServer(t)
	setupMongoDB(t)

	defer func() {
		if err := mongoapi.Drop("policyData.ues.eirData"); err != nil {
			panic(err)
		}
	}()

	filter := bson.M{"pei": nil}
	pei1 := bson.M{
		"pei": "imei-42", "supi": "imsi-208930000000042", "gpsi": "msisdn-00042",
		"equipment_status": "BLACKLISTED",
	}
	pei2 := bson.M{
		"pei": "imei-012345678901234", "supi": "imsi-012345678901234", "gpsi": "msisdn-12345",
		"equipment_status": "WHITELISTED",
	}
	pei3 := bson.M{
		"pei": "imei-43", "supi": "imsi-208930000000043", "gpsi": "msisdn-00043",
		"equipment_status": "BLACKLISTED",
	}
	filters := []bson.M{filter, filter, filter}
	peis := []map[string]interface{}{pei1, pei2, pei3}
	err := mongoapi.RestfulAPIPutMany("policyData.ues.eirData", filters, peis)
	assert.Nil(t, err)

	reqUri := factory.EirDrResUriPrefix +
		"/equipment-status?pei=imei-012345678901234&supi=imsi-012345678901234&gpsi=msisdn-12345"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	expected_message := util.ToBsonM(eir_api_service.EIREquipmentStatusGetResponse{
		Status: "WHITELISTED",
	})
	t.Run("EquipmentStatus", func(t *testing.T) {
		json_message := eir_api_service.EIREquipmentStatusGetResponse{}

		err := json.Unmarshal(rsp.Body.Bytes(), &json_message)
		assert.Nil(t, err)

		message := util.ToBsonM(json_message)

		require.Equal(t, expected_message, message)
		require.Equal(t, http.StatusOK, rsp.Code)
	})
}

func TestEIR_EquipmentStatus_NotFoundEquipmentStatus(t *testing.T) {
	server := setupHttpServer(t)
	setupMongoDB(t)

	defer func() {
		if err := mongoapi.Drop("policyData.ues.eirData"); err != nil {
			panic(err)
		}
	}()

	reqUri := factory.EirDrResUriPrefix + "/equipment-status?pei=imei-012345678901234"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	expected_message := util.ToBsonM(models.ProblemDetails{
		Title:  "The equipment identify checking has failed",
		Status: http.StatusNotFound,
		Detail: "The Equipment Status wasn't found",
		Cause:  "ERROR_EQUIPMENT_UNKNOWN",
	})
	t.Run("EquipmentStatus", func(t *testing.T) {
		json_message := models.ProblemDetails{}

		err := json.Unmarshal(rsp.Body.Bytes(), &json_message)
		assert.Nil(t, err)

		message := util.ToBsonM(json_message)

		require.Equal(t, expected_message, message)
		require.Equal(t, http.StatusNotFound, rsp.Code)
	})
}

func TestEIR_EquipmentStatus_MissingPEI(t *testing.T) {
	server := setupHttpServer(t)
	setupMongoDB(t)

	defer func() {
		if err := mongoapi.Drop("policyData.ues.eirData"); err != nil {
			panic(err)
		}
	}()

	reqUri := factory.EirDrResUriPrefix + "/equipment-status"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	expected_message := util.ToBsonM(models.ProblemDetails{
		Title:  "The equipment identify checking has failed",
		Status: http.StatusBadRequest,
		Detail: "The PEI is missing",
		Cause:  "MANDATORY_IE_MISSING",
		InvalidParams: []models.InvalidParam{{
			Param:  "PEI",
			Reason: "The PEI is missing",
		}},
	})
	t.Run("EquipmentStatus", func(t *testing.T) {
		json_message := models.ProblemDetails{}

		err := json.Unmarshal(rsp.Body.Bytes(), &json_message)
		assert.Nil(t, err)

		message := util.ToBsonM(json_message)

		require.Equal(t, expected_message, message)
		require.Equal(t, http.StatusNotFound, rsp.Code)
	})
}

func TestEIR_EquipmentStatus_NotFoundEquipmentStatus_WithDefaultStatusBlack(t *testing.T) {
	server := setupHttpServerWithDefaultStatus(t, "BLACKLISTED")
	setupMongoDB(t)

	defer func() {
		if err := mongoapi.Drop("policyData.ues.eirData"); err != nil {
			panic(err)
		}
	}()

	reqUri := factory.EirDrResUriPrefix + "/equipment-status?pei=imei-012345678901234"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	expected_message := util.ToBsonM(eir_api_service.EIREquipmentStatusGetResponse{
		Status: "BLACKLISTED",
	})
	t.Run("EquipmentStatus", func(t *testing.T) {
		json_message := eir_api_service.EIREquipmentStatusGetResponse{}

		err := json.Unmarshal(rsp.Body.Bytes(), &json_message)
		assert.Nil(t, err)

		message := util.ToBsonM(json_message)

		require.Equal(t, expected_message, message)
		require.Equal(t, http.StatusOK, rsp.Code)
	})
}

func TestEIR_EquipmentStatus_NotFoundEquipmentStatus_WithDefaultStatusWhite(t *testing.T) {
	server := setupHttpServerWithDefaultStatus(t, "WHITELISTED")
	setupMongoDB(t)

	defer func() {
		if err := mongoapi.Drop("policyData.ues.eirData"); err != nil {
			panic(err)
		}
	}()

	reqUri := factory.EirDrResUriPrefix + "/equipment-status?pei=imei-012345678901234"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	expected_message := util.ToBsonM(eir_api_service.EIREquipmentStatusGetResponse{
		Status: "WHITELISTED",
	})
	t.Run("EquipmentStatus", func(t *testing.T) {
		json_message := eir_api_service.EIREquipmentStatusGetResponse{}

		err := json.Unmarshal(rsp.Body.Bytes(), &json_message)
		assert.Nil(t, err)

		message := util.ToBsonM(json_message)

		require.Equal(t, expected_message, message)
		require.Equal(t, http.StatusOK, rsp.Code)
	})
}

func TestEIR_EquipmentStatus_URITooLong(t *testing.T) {
	server := setupHttpServer(t)
	setupMongoDB(t)

	defer func() {
		if err := mongoapi.Drop("policyData.ues.eirData"); err != nil {
			panic(err)
		}
	}()

	imei := "imei-" + strings.Repeat("0", 4096)
	reqUri := factory.EirDrResUriPrefix + "/equipment-status?pei=" + imei

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	expected_message := util.ToBsonM(models.ProblemDetails{
		Title:  "The equipment identify checking has failed",
		Status: http.StatusRequestURITooLong,
		Detail: "URI Too Long",
		Cause:  "INCORRECT_URI_LENGTH",
	})
	t.Run("EquipmentStatus", func(t *testing.T) {
		json_message := models.ProblemDetails{}

		err := json.Unmarshal(rsp.Body.Bytes(), &json_message)
		assert.Nil(t, err)

		message := util.ToBsonM(json_message)

		require.Equal(t, expected_message, message)
		require.Equal(t, http.StatusRequestURITooLong, rsp.Code)
	})
}

func TestEIR_EquipmentStatus_WithoutDatabase(t *testing.T) {
	server := setupHttpServer(t)
	err := mongoapi.Client.Disconnect(context.Background()) // The reason of the error
	if err != nil {
		logger.UtilLog.Errorf("Failed to properly disconnect from the database")
	}

	reqUri := factory.EirDrResUriPrefix + "/equipment-status?pei=012345678901234"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	expected_message := util.ToBsonM(models.ProblemDetails{
		Title:  "The equipment identify checking has failed",
		Status: http.StatusInternalServerError,
		Detail: "RestfulAPIGetOne err: client is disconnected",
		Cause:  "SYSTEM_FAILURE",
	})
	t.Run("EquipmentStatus", func(t *testing.T) {
		json_message := models.ProblemDetails{}

		err := json.Unmarshal(rsp.Body.Bytes(), &json_message)
		assert.Nil(t, err)

		message := util.ToBsonM(json_message)

		logger.UtilLog.Info(message)
		logger.UtilLog.Info(rsp.Code)
		require.Equal(t, expected_message, message)
		require.Equal(t, http.StatusInternalServerError, rsp.Code)
	})
}
