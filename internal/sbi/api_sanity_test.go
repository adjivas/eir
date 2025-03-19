package sbi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	db "github.com/adjivas/eir/internal/database"
	"github.com/adjivas/eir/internal/logger"
	"github.com/adjivas/eir/pkg/factory"
	util_logger "github.com/free5gc/util/logger"
	"github.com/free5gc/util/mongoapi"
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

func TestEIR_EquipementStatus(t *testing.T) {
	server := setupHttpServer(t)
	reqUri := factory.EirDrResUriPrefix + "/equipement-status"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	t.Run("EquipementStatus", func(t *testing.T) {
		require.Equal(t, http.StatusNotImplemented, rsp.Code)
	})
}
