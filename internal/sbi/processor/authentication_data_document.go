/*
 * Nudr_DataRepository API OpenAPI file
 *
 * Unified Data Repository Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package processor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/udr/internal/logger"
	"github.com/free5gc/udr/internal/util"
)

func (p *Processor) ModifyAuthenticationProcedure(
	c *gin.Context, collName string, ueId string, patchItem []models.PatchItem,
) {
	filter := bson.M{"ueId": ueId}
	if err := patchDataToDBAndNotify(collName, ueId, patchItem, filter); err != nil {
		logger.DataRepoLog.Errorf("ModifyAuthenticationProcedure err: %+v", err)
		c.JSON(http.StatusInternalServerError, util.ProblemDetailsModifyNotAllowed(""))
	}
	c.Status(http.StatusNoContent)
}

func (p *Processor) QueryAuthSubsDataProcedure(c *gin.Context, collName string, ueId string) {
	filter := bson.M{"ueId": ueId}
	data, pd := getDataFromDB(collName, filter)
	if pd != nil {
		if pd.Status == http.StatusNotFound {
			logger.DataRepoLog.Warnf("QueryAuthSubsDataProcedure err: %s", pd.Title)
		} else {
			logger.DataRepoLog.Errorf("QueryAuthSubsDataProcedure err: %s", pd.Detail)
		}
		c.JSON(int(pd.Status), pd)
		return
	}
	c.JSON(http.StatusOK, data)
}
