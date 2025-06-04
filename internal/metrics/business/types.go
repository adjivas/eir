package business

// Global metric information
const (
	SUBSYSTEM_NAME = "business"
)

// Collectors information
const (
	FAIL_EQUIPMENT_STATUS_COUNTER_NAME = "fail_equipment_status_total"
	FAIL_EQUIPMENT_STATUS_COUNTER_DESC = "Counter of equipment status failure"

	SUCCESS_EQUIPMENT_STATUS_COUNTER_NAME = "success_equipment_status_total"
	SUCCESS_EQUIPMENT_STATUS_COUNTER_DESC = "success Counter of equipment status success"
)

// Label names
const (
	EIR_STATUS_LABEL = "status"
	EIR_TYPE_LABEL   = "type"
)

// Metrics Values
const (
	EIR_ERROR = "error"
	EIR_WARN  = "warn"
)

// Potential Causes
const (
	PEI_NOT_FOUND     = "pei not found"
	DB_SYSTEM_FAILURE = "system failure"
	DB_UNSPECIFIED    = "unspecified"
)
