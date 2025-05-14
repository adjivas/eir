package sbi

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	OutboundReqCounterName = "outbound_request_total"
	OutboundReqCounterDesc = "Total number of SBI outbound requests attempted or sent by the AMF"
)

const (
	SUBSYSTEM_NAME = "sbi"
)

var (
	OutboundReqCounter prometheus.CounterVec
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
