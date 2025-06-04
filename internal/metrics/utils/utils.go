package utils

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
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

func getValueFromGauge(c prometheus.Gauge) (float64, error) {
	m := &dto.Metric{}
	if err := c.Write(m); err != nil {
		return 0, err
	}
	return m.GetGauge().GetValue(), nil
}

func GetGaugeVecValue(gaugeName string, gauge *prometheus.GaugeVec, labels prometheus.Labels) (float64, error) {
	foundGauge, err := gauge.GetMetricWith(labels)

	if err != nil {
		return 0, fmt.Errorf("could not retrieve the %s gauge", gaugeName)
	}

	gaugeValue, err := getValueFromGauge(foundGauge)

	if err != nil {
		return 0, fmt.Errorf("failed to get %s gauge value", gaugeName)
	}

	return gaugeValue, nil
}

func FormatStatus(statusCode int) string {
	code := http.StatusInternalServerError
	if statusCode != 0 {
		code = statusCode
	}

	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}

// readStringPtr return the value of the string pointer if non-nil. Returns an empty string otherwise
func ReadStringPtr(strPtr *string) string {
	if strPtr == nil {
		temp := ""
		strPtr = &temp
	}
	return *strPtr
}
