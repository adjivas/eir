package processor

import (
	"net/http"

	"github.com/adjivas/eir/internal/logger"
	business_metrics "github.com/adjivas/eir/internal/metrics/business"
	"github.com/adjivas/eir/internal/metrics/sbi"
	"github.com/adjivas/eir/internal/util"
	eir_api_service "github.com/free5gc/openapi/eir/EIRService"
	"github.com/free5gc/openapi/models"
	"github.com/gin-gonic/gin"
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
		business_metrics.IncrEquipmentStatusSuccessCounter()
		c.JSON(http.StatusOK, response)
		return
	}
	switch err_database.Cause {
	case "DATA_NOT_FOUND":
		if defaultStatus := p.App.Config().Configuration.DefaultStatus; defaultStatus != "" {
			logger.ProcLog.Warnf("The Equipment Status wasn't found, the default %s is returned", defaultStatus)
			response := util.ToBsonM(eir_api_service.EIREquipmentStatusGetResponse{
				Status: defaultStatus,
			})
			business_metrics.IncrEquipmentStatusFailCounter(business_metrics.EIR_WARN, business_metrics.PEI_NOT_FOUND)
			c.JSON(http.StatusOK, response)
		} else {
			problemDetail := models.ProblemDetails{
				Title:  "The equipment identify checking has failed",
				Status: http.StatusNotFound,
				Detail: "The Equipment Status wasn't found",
				Cause:  "ERROR_EQUIPMENT_UNKNOWN",
			}
			business_metrics.IncrEquipmentStatusFailCounter(business_metrics.EIR_ERROR, business_metrics.PEI_NOT_FOUND)
			c.Set(sbi.IN_PB_DETAILS_CTX_STR, problemDetail.Cause)
			c.JSON(http.StatusNotFound, problemDetail)
		}
	case "SYSTEM_FAILURE":
		logger.ProcLog.Errorf("The database has failed with [%v]", err_database.Detail)
		problemDetail := models.ProblemDetails{
			Title:  "The equipment identify checking has failed",
			Status: http.StatusInternalServerError,
			Detail: err_database.Detail,
			Cause:  "SYSTEM_FAILURE",
		}
		business_metrics.IncrEquipmentStatusFailCounter(business_metrics.EIR_ERROR, business_metrics.DB_SYSTEM_FAILURE)
		c.JSON(http.StatusInternalServerError, problemDetail)
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, problemDetail.Cause)
	default:
		logger.ProcLog.Errorf("The NF has a unspecified failure with [%+v]", err_database)
		problemDetail := models.ProblemDetails{
			Title:  "The equipment identify checking has failed",
			Status: http.StatusInternalServerError,
			Detail: err_database.Detail,
			Cause:  "INSUFFICIENT_RESOURCES",
		}
		business_metrics.IncrEquipmentStatusFailCounter(business_metrics.EIR_ERROR, business_metrics.DB_UNSPECIFIED)
		c.JSON(http.StatusInternalServerError, problemDetail)
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, problemDetail.Cause)
	}
}
