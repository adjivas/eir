package processor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/adjivas/eir/internal/logger"
	"github.com/free5gc/openapi/models"
)

func (p *Processor) GetEirEquipementStatusProcedure(c *gin.Context, collName string, pei string) {
	filter := bson.M{"pei": pei}
	data, p_equipement_status := p.GetDataFromDB(collName, filter)
	if p_equipement_status != nil {
		problemDetail := models.ProblemDetails{
			Title: "The equipment identify checking has failed",
			Status: http.StatusNotFound,
			Detail: "The Equipment Status wasn't found",
			Cause:  "ERROR_EQUIPMENT_UNKNOWN",
		}
		logger.CallbackLog.Errorf("The PEI is missing")
		c.JSON(http.StatusNotFound, problemDetail)
	} else {
		c.JSON(http.StatusOK, data)
	}
}
