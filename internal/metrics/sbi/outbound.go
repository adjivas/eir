package sbi

import (
	"github.com/adjivas/eir/internal/metrics/utils"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

func GetSbiOutboundMetrics(namespace string) []prometheus.Collector {
	var metrics []prometheus.Collector

	OutboundReqCounter = *prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      OutboundReqCounterName,
			Help:      OutboundReqCounterDesc,
		},
		[]string{OUT_TARGET_SERVICE_NAME_LABEL, OUT_STATUS_CODE_LABEL, OUT_METHOD_LABEL},
	)

	metrics = append(metrics, OutboundReqCounter)

	return metrics
}

func IncrOutboundReqCounter(metricInfo *OutboundMetricBasicInfo) {

	status := ""
	if metricInfo.StatusCode != 0 {
		status = utils.FormatStatus(metricInfo.StatusCode)
	} else {
		status = utils.FormatStatus(http.StatusInternalServerError)
	}

	OutboundReqCounter.With(prometheus.Labels{
		OUT_TARGET_SERVICE_NAME_LABEL: metricInfo.TargetServiceName,
		OUT_STATUS_CODE_LABEL:         status,
		OUT_METHOD_LABEL:              metricInfo.Method,
	}).Add(1)
}

func SbiMetricHook(method string, serviceName string, status int) {

	info := OutboundMetricBasicInfo{
		TargetServiceName: serviceName,
		StatusCode:        status,
		Method:            method,
	}

	IncrOutboundReqCounter(&info)
}
