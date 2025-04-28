package util

import (
	"net/http"

	"github.com/free5gc/openapi/models"
)

const (
	EQUIPMENT_UNKNOWN = "Data not found"
)

func ProblemDetailsNotFound(cause string) *models.ProblemDetails {
	title := EQUIPMENT_UNKNOWN
	return &models.ProblemDetails{
		Title:  title,
		Status: http.StatusNotFound,
		Cause:  cause,
	}
}
