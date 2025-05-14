package sbi

type MetricReporter interface {
	Report(*OutboundMetricBasicInfo)
}
type defaultReporter struct{}

func (defaultReporter) Report(info *OutboundMetricBasicInfo) {
	IncrOutboundReqCounter(info)
}

func NewDefaultReporter() MetricReporter {
	return defaultReporter{}
}
