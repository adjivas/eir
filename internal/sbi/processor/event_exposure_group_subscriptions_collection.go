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
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi/models"
	udr_context "github.com/free5gc/udr/internal/context"
	"github.com/free5gc/udr/internal/util"
	"github.com/free5gc/udr/pkg/factory"
)

func (p *Processor) CreateEeGroupSubscriptionsProcedure(
	c *gin.Context, ueGroupId string, EeSubscription models.EeSubscription,
) {
	udrSelf := udr_context.GetSelf()

	value, ok := udrSelf.UEGroupCollection.Load(ueGroupId)
	if !ok {
		udrSelf.UEGroupCollection.Store(ueGroupId, new(udr_context.UEGroupSubsData))
		value, _ = udrSelf.UEGroupCollection.Load(ueGroupId)
	}
	UEGroupSubsData := value.(*udr_context.UEGroupSubsData)
	if UEGroupSubsData.EeSubscriptions == nil {
		UEGroupSubsData.EeSubscriptions = make(map[string]*models.EeSubscription)
	}

	newSubscriptionID := strconv.Itoa(udrSelf.EeSubscriptionIDGenerator)
	UEGroupSubsData.EeSubscriptions[newSubscriptionID] = &EeSubscription
	udrSelf.EeSubscriptionIDGenerator++

	/* Contains the URI of the newly created resource, according
	   to the structure: {apiRoot}/nudr-dr/v1/subscription-data/group-data/{ueGroupId}/ee-subscriptions */
	locationHeader := fmt.Sprintf("%s"+factory.UdrDrResUriPrefix+"/subscription-data/group-data/%s/ee-subscriptions/%s",
		udrSelf.GetIPv4GroupUri(udr_context.NUDR_DR), ueGroupId, newSubscriptionID)

	c.Header("Location", locationHeader)
	c.JSON(http.StatusCreated, EeSubscription)
}

func (p *Processor) QueryEeGroupSubscriptionsProcedure(c *gin.Context, ueGroupId string) {
	udrSelf := udr_context.GetSelf()

	value, ok := udrSelf.UEGroupCollection.Load(ueGroupId)
	if !ok {
		pd := util.ProblemDetailsNotFound("USER_NOT_FOUND")
		c.JSON(int(pd.Status), pd)
		return
	}

	UEGroupSubsData := value.(*udr_context.UEGroupSubsData)
	var eeSubscriptionSlice []models.EeSubscription

	for _, v := range UEGroupSubsData.EeSubscriptions {
		eeSubscriptionSlice = append(eeSubscriptionSlice, *v)
	}

	if len(eeSubscriptionSlice) == 0 {
		pd := util.ProblemDetailsUpspecified("")
		c.JSON(int(pd.Status), pd)
		return
	}
	c.JSON(http.StatusOK, eeSubscriptionSlice)
}
