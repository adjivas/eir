package processor

import (
	"net/http"

	"github.com/adjivas/eir/internal/logger"
	eir_models "github.com/adjivas/eir/internal/models"
	"github.com/adjivas/eir/internal/util"
	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi/models"
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
		problemDetail := models.ProblemDetails{
			Title:  "The equipment identify checking has failed",
			Status: http.StatusNotFound,
			Detail: "The Equipment Status wasn't found",
			Cause:  "ERROR_EQUIPMENT_UNKNOWN",
		}
		logger.CallbackLog.Errorf("The Equipment Status wasn't found")
		c.JSON(http.StatusNotFound, problemDetail)
	} else {
		response := util.ToBsonM(eir_models.EirResponseData{
			Status: data["equipement_status"].(string),
		})
		c.JSON(http.StatusOK, response)
	}
}
