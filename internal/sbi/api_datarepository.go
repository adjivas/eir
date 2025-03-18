/*
 * N5g_DataRepository API OpenAPI file
 *
 * Unified Data Repository Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (s://openapi-generator.tech)
 */

package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	// eir_context "github.com/adjivas/eir/internal/context"
	"github.com/adjivas/eir/internal/logger"
	"github.com/adjivas/eir/internal/util"
)

func (s *Server) getDataRepositoryRoutes() []Route {
	return []Route{
		{
			"Index",
			"GET",
			"/",
			Index,
		},
	}
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Hello World!")
}

// HTTPAmfContext3gpp - To modify the AMF context data of a UE using 3gpp access in the EIR
func (s *Server) HandleAmfContext3gpp(c *gin.Context) {
	var patchItemArray []models.PatchItem

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.DataRepoLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&patchItemArray, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.DataRepoLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	logger.DataRepoLog.Tracef("Handle AmfContext3gpp")
	ueId := c.Params.ByName("ueId")
	if ueId == "" {
		util.EmptyUeIdProblemJson(c)
		return
	}
}
