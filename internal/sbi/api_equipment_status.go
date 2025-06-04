package sbi

import (
	"net/http"

	"github.com/adjivas/eir/internal/logger"
	"github.com/free5gc/openapi/models"
	"github.com/gin-gonic/gin"
)

const maxURILength = 1024

func (s *Server) getEquipmentStatusRoutes() []Route {
	return []Route{
		{
			"Index",
			"GET",
			"/",
			Index,
		},
		{
			"EquipmentStatus",
			"GET",
			"/equipment-status",
			s.HandleQueryEirEquipmentStatus,
		},
	}
}

func URILengthLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if size := len(c.Request.URL.String()); size > maxURILength {
			problemDetail := models.ProblemDetails{
				Title:  "The equipment identify checking has failed",
				Status: http.StatusRequestURITooLong,
				Detail: "URI Too Long",
				Cause:  "INCORRECT_URI_LENGTH",
			}
			logger.HttpLog.Errorf("The Request URI is too long (%d>%d)", size, maxURILength)
			c.JSON(http.StatusRequestURITooLong, problemDetail)
			c.Abort()
			return
		}
		c.Next()
	}
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Hello World!")
}

func (s *Server) HandleQueryEirEquipmentStatus(c *gin.Context) {
	logger.EquipmentStatusLog.Tracef("Handle EirEquipmentStatus")

	collName := "policyData.ues.eirData"
	pei := c.Query("pei")
	supi := c.DefaultQuery("supi", "")
	gpsi := c.DefaultQuery("gpsi", "")
	if pei == "" {
		problemDetail := models.ProblemDetails{
			Title:  "The equipment identify checking has failed",
			Status: http.StatusBadRequest,
			Detail: "The PEI is missing",
			Cause:  "MANDATORY_IE_MISSING",
			InvalidParams: []models.InvalidParam{{
				Param:  "PEI",
				Reason: "The PEI is missing",
			}},
		}
		logger.HttpLog.Errorf("The PEI is missing")
		c.JSON(http.StatusNotFound, problemDetail)
	} else {
		s.eir.Processor().GetEirEquipmentStatusProcedure(c, collName, pei, supi, gpsi)
	}
}
