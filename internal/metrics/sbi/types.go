package sbi

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	InboundReqCounterName = "inbound_request_total"
	InboundReqCounterDesc = "Total number of SBI inbound requests received by the AMF"
	InboundReqHistogramName = "inbound_request_duration_seconds"
	InboundReqHistogramDesc = "Histogram of request latencies"

	OutboundReqCounterName = "outbound_request_total"
	OutboundReqCounterDesc = "Total number of SBI outbound requests attempted or sent by the AMF"
)

const (
	SUBSYSTEM_NAME = "sbi"
)

var (
	InboundReqCounter      prometheus.CounterVec
	InboundRequestDuration prometheus.HistogramVec

	OutboundReqCounter prometheus.CounterVec
)





// Labels names for the inbound sbi metrics


const (
	IN_STATUS_CODE_LABEL   = "status_code"
	IN_METHOD_LABEL        = "method"
	IN_REQUESTED_URL_LABEL = "requested_url"
	IN_CAUSE_LABEL         = "cause"
	IN_PATH_LABEL          = "path"
	IN_PB_DETAILS_CTX_STR  = "problem"
)

// Labels names for the outbound sbi metrics
const (
	OUT_TARGET_SERVICE_NAME_LABEL = "target_service_name"
	OUT_STATUS_CODE_LABEL         = "status_code"
	OUT_METHOD_LABEL              = "method"
)

type OutboundMetricBasicInfo struct {
	StatusCode        int    `json:"status_code"`
	TargetServiceName string `json:"target_service_name"`
	Method            string `json:"method"`
}
