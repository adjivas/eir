// Package metrics sets and initializes Prometheus metrics.
package metrics

import (
	"github.com/adjivas/eir/internal/metrics/business"
	"github.com/adjivas/eir/internal/metrics/sbi"
	"github.com/adjivas/eir/pkg/factory"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// Init initializes all Prometheus metrics
func Init(cfg *factory.Config) *prometheus.Registry {
	reg := prometheus.NewRegistry()

	namespace := cfg.GetMetricsNamespace()

	globalLabels := prometheus.Labels{
		NF_TYPE_LABEL: NF_TYPE_VALUE,
	}

	wrappedReg := prometheus.WrapRegistererWith(globalLabels, reg)

	wrappedReg.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	var eirMetrics []prometheus.Collector

	// Append here the collector you want to register to the prometheus registry
	eirMetrics = append(eirMetrics, sbi.GetSbiOutboundMetrics(namespace)...)
	eirMetrics = append(eirMetrics, sbi.GetSbiInboundMetrics(namespace)...)

	eirMetrics = append(eirMetrics, business.GetEquipmentStatusMetrics(namespace)...)

	initMetric(eirMetrics, wrappedReg)

	return reg
}

func initMetric(metrics []prometheus.Collector, reg prometheus.Registerer) {
	for _, metric := range metrics {
		reg.MustRegister(metric)
	}
}
