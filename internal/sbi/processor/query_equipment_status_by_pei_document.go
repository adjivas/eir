package processor

import (
	"net/http"

	"github.com/adjivas/eir/internal/logger"
	"github.com/adjivas/eir/internal/util"
	"github.com/gin-gonic/gin"

	eir_api_service "github.com/free5gc/openapi/eir/EIRService"
	"github.com/free5gc/openapi/models"
)

func (p *Processor) GetEirEquipmentStatusProcedure(c *gin.Context, collName string,
	pei string, supi string, gpsi string,
) {
	filter := map[string]interface{}{
		"pei": pei,
	}
	if supi != "" {
		filter["supi"] = supi
	}
	if gpsi != "" {
		filter["gpsi"] = gpsi
	}

	data, err_database := p.DbConnector.GetDataFromDB(collName, filter)
	if err_database == nil {
		response := util.ToBsonM(eir_api_service.EIREquipmentStatusGetResponse{
			Status: data["equipment_status"].(string),
		})
		c.JSON(http.StatusOK, response)
	} else {
		switch err_database.Cause {
		case "DATA_NOT_FOUND":
			if defaultStatus := p.App.Config().Configuration.DefaultStatus; defaultStatus != "" {
				logger.ProcLog.Warnf("The Equipment Status wasn't found, the default %s is returned", defaultStatus)
				response := util.ToBsonM(eir_api_service.EIREquipmentStatusGetResponse{
					Status: defaultStatus,
				})
				c.JSON(http.StatusOK, response)
			} else {
				logger.ProcLog.Errorln("The Equipment Status wasn't found")
				problemDetail := models.ProblemDetails{
					Title:  "The equipment identify checking has failed",
					Status: http.StatusNotFound,
					Detail: "The Equipment Status wasn't found",
					Cause:  "ERROR_EQUIPMENT_UNKNOWN",
				}
				c.JSON(http.StatusNotFound, problemDetail)
			}
		case "SYSTEM_FAILURE":
			logger.ProcLog.Errorf("The database has failed with [%v]", err_database.Detail)
			problemDetail := models.ProblemDetails{
				Title:  "The equipment identify checking has failed",
				Status: http.StatusInternalServerError,
				Detail: err_database.Detail,
				Cause:  "INSUFFICIENT_RESOURCES",
			}
			c.JSON(http.StatusInternalServerError, problemDetail)
		default:
			logger.ProcLog.Errorf("The NF has a unspecified failure with [%+v]", err_database)
			problemDetail := models.ProblemDetails{
				Title:  "The equipment identify checking has failed",
				Status: http.StatusInternalServerError,
				Detail: err_database.Detail,
				Cause:  "INSUFFICIENT_RESOURCES",
			}
			c.JSON(http.StatusInternalServerError, problemDetail)
		}
	}
}
