package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adjivas/eir/internal/logger"
	"github.com/free5gc/openapi/models"
)

func (s *Server) getEquipementStatusRoutes() []Route {
	return []Route{
		{
			"Index",
			"GET",
			"/",
			Index,
		},
		{
			"EquipementStatus",
			"GET",
			"/equipement-status",
			s.HandleQueryEirEquipementStatus,
		},
	}
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Hello World!")
}

func (s *Server) HandleQueryEirEquipementStatus(c *gin.Context) {
	logger.EquipementStatusLog.Tracef("Handle EirEquipementStatus")

	collName := "policyData.ues.eirData"
	pei := c.Query("pei")
	supi := c.DefaultQuery("supi", "")
	gpsi := c.DefaultQuery("gpsi", "")
	if pei == "" {
		problemDetail := models.ProblemDetails{
			Title: "The equipment identify checking has failed",
			Status: http.StatusNotFound,
			Detail: "The PEI is missing",
			Cause:  "ERROR_EQUIPMENT_UNKNOWN",
		}
		logger.CallbackLog.Errorf("The PEI is missing")
		c.JSON(http.StatusNotFound, problemDetail)
	} else {
		s.Processor().GetEirEquipementStatusProcedure(c, collName, pei, supi, gpsi)
    }
}
