package sbi

// Outbound
import (
	"github.com/adjivas/eir/internal/metrics/utils"
	"github.com/prometheus/client_golang/prometheus"
)

func GetSbiOutboundMetrics(namespace string) []prometheus.Collector {
	var metrics []prometheus.Collector

	OutboundReqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      OUT_BOUND_REQ_COUNTER_NAME,
			Help:      OUT_BOUND_REQ_COUNTER_DESC,
		},
		[]string{OUT_TARGET_SERVICE_NAME_LABEL, OUT_STATUS_CODE_LABEL, OUT_METHOD_LABEL},
	)

	metrics = append(metrics, OutboundReqCounter)

	OutboundRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      OUT_BOUND_REQ_HISTOGRAM_NAME,
			Help:      OUT_BOUND_REQ_HISTOGRAM_DESC,
			Buckets: []float64{
				0.0001,
				0.0050,
				0.0100,
				0.0200,
				0.0250,
				0.0500,
			},
		},
		[]string{OUT_METHOD_LABEL, OUT_TARGET_SERVICE_NAME_LABEL, OUT_STATUS_CODE_LABEL},
	)

	metrics = append(metrics, OutboundRequestDuration)

	return metrics
}

func IncrOutboundReqCounter(method string, serviceName string, statusCode int) {

	OutboundReqCounter.With(prometheus.Labels{
		OUT_TARGET_SERVICE_NAME_LABEL: serviceName,
		OUT_STATUS_CODE_LABEL:         utils.FormatStatus(statusCode),
		OUT_METHOD_LABEL:              method,
	}).Add(1)
}

func IncrOutboundReqDurationCounter(method string, serviceName string, statusCode int, duration float64) {
	OutboundRequestDuration.With(prometheus.Labels{
		OUT_TARGET_SERVICE_NAME_LABEL: serviceName,
		OUT_METHOD_LABEL:              method,
		OUT_STATUS_CODE_LABEL:         utils.FormatStatus(statusCode),
	}).Observe(duration)
}

func SbiMetricHook(method string, serviceName string, statusCode int, duration float64) {
	IncrOutboundReqCounter(method, serviceName, statusCode)
	IncrOutboundReqDurationCounter(method, serviceName, statusCode, duration)
}
