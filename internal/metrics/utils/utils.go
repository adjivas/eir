package utils

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"net/http"
)

const (
	SuccessMetric = "successful"
	FailureMetric = "failure"
)

func getValueFromCounter(c prometheus.Counter) (float64, error) {
	m := &dto.Metric{}
	if err := c.Write(m); err != nil {
		return 0, err
	}
	return m.GetCounter().GetValue(), nil
}

func GetCounterVecValue(counterName string, counter prometheus.CounterVec, labels prometheus.Labels) (float64, error) {
	foundCounter, err := counter.GetMetricWith(labels)

	if err != nil {
		return 0, fmt.Errorf("could not retrieve the %s counter", counterName)
	}

	counterValue, err := getValueFromCounter(foundCounter)

	if err != nil {
		return 0, fmt.Errorf("failed to get %s counter value", counterName)
	}

	return counterValue, nil
}

func FormatStatus(statusCode int) string {
	return fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
}
