package business

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// EirEquipmentStatusFailCounter Counter of failure,
	// labeled by status_label type_label
	EirEquipmentStatusFailCounter *prometheus.CounterVec

	// EirEquipmentStatusSuccessCounter Counter of success,
	// labeled by status_label type_label
	EirEquipmentStatusSuccessCounter prometheus.Counter
)

func GetEquipmentStatusMetrics(namespace string) []prometheus.Collector {
	var collectors []prometheus.Collector

	EirEquipmentStatusFailCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      FAIL_EQUIPMENT_STATUS_COUNTER_NAME,
			Help:      FAIL_EQUIPMENT_STATUS_COUNTER_DESC,
		},
		[]string{EIR_STATUS_LABEL, EIR_TYPE_LABEL},
	)

	collectors = append(collectors, EirEquipmentStatusFailCounter)

	EirEquipmentStatusSuccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: SUCCESS_EQUIPMENT_STATUS_COUNTER_NAME,
			Help: SUCCESS_EQUIPMENT_STATUS_COUNTER_DESC,
		},
	)

	collectors = append(collectors, EirEquipmentStatusSuccessCounter)

	return collectors
}

func IncrEquipmentStatusFailCounter(eirStatus string, eirType string) {
	EirEquipmentStatusFailCounter.With(prometheus.Labels{
		EIR_STATUS_LABEL: eirStatus,
		EIR_TYPE_LABEL:   eirType,
	}).Add(1)
}

func IncrEquipmentStatusSuccessCounter() {
	EirEquipmentStatusSuccessCounter.Inc()
}
