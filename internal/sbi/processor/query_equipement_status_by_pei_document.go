package processor

import (
	"net/http"

	"github.com/adjivas/eir/internal/logger"
	eir_models "github.com/adjivas/eir/internal/models"
	"github.com/adjivas/eir/internal/util"
	"github.com/free5gc/openapi/models"
	"github.com/gin-gonic/gin"
)

func (p *Processor) GetEirEquipementStatusProcedure(c *gin.Context, collName string,
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

	data, p_equipement_status := p.GetDataFromDB(collName, filter)
	if p_equipement_status != nil {

		if defaultStatus := p.Config().Configuration.DefaultStatus; defaultStatus != "" {
			logger.CallbackLog.Warnf("The Equipment Status wasn't found, the default %s is returned", defaultStatus)
			response := util.ToBsonM(eir_models.EirResponseData{
				Status: defaultStatus,
			})
			c.JSON(http.StatusOK, response)
		} else {
			logger.CallbackLog.Errorln("The Equipment Status wasn't found")
			problemDetail := models.ProblemDetails{
				Title:  "The equipment identify checking has failed",
				Status: http.StatusNotFound,
				Detail: "The Equipment Status wasn't found",
				Cause:  "ERROR_EQUIPMENT_UNKNOWN",
			}
			c.JSON(http.StatusNotFound, problemDetail)
		}
	} else {
		response := util.ToBsonM(eir_models.EirResponseData{
			Status: data["equipement_status"].(string),
		})
		c.JSON(http.StatusOK, response)
	}
}
