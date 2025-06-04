package util

import (
	"net/http"

	eir_context "github.com/adjivas/eir/internal/context"
	"github.com/adjivas/eir/internal/logger"
	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi/models"
)

type RouterAuthorizationCheck struct {
	serviceName models.ServiceName
}

func NewRouterAuthorizationCheck(serviceName models.ServiceName) *RouterAuthorizationCheck {
	return &RouterAuthorizationCheck{
		serviceName: serviceName,
	}
}

func (rac *RouterAuthorizationCheck) Check(c *gin.Context, eirContext eir_context.NFContext) {
	token := c.Request.Header.Get("Authorization")
	err := eirContext.AuthorizationCheck(token, rac.serviceName)
	if err != nil {
		logger.UtilLog.Debugf("RouterAuthorizationCheck: Check Unauthorized: %s", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	logger.UtilLog.Debugf("RouterAuthorizationCheck: Check Authorized")
}
