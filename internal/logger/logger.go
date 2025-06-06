package logger

import (
	logger_util "github.com/free5gc/util/logger"
	"github.com/sirupsen/logrus"
)

var (
	Log                *logrus.Logger
	NfLog              *logrus.Entry
	MainLog            *logrus.Entry
	InitLog            *logrus.Entry
	CfgLog             *logrus.Entry
	EquipmentStatusLog *logrus.Entry
	UtilLog            *logrus.Entry
	HttpLog            *logrus.Entry
	ConsumerLog        *logrus.Entry
	GinLog             *logrus.Entry
	ProcLog            *logrus.Entry
	SBILog             *logrus.Entry
	DbLog              *logrus.Entry
)

func init() {
	fieldsOrder := []string{
		logger_util.FieldNF,
		logger_util.FieldCategory,
	}

	Log = logger_util.New(fieldsOrder)
	NfLog = Log.WithField(logger_util.FieldNF, "EIR")
	MainLog = NfLog.WithField(logger_util.FieldCategory, "Main")
	InitLog = NfLog.WithField(logger_util.FieldCategory, "Init")
	CfgLog = NfLog.WithField(logger_util.FieldCategory, "CFG")
	GinLog = NfLog.WithField(logger_util.FieldCategory, "GIN")
	ConsumerLog = NfLog.WithField(logger_util.FieldCategory, "Consumer")
	EquipmentStatusLog = NfLog.WithField(logger_util.FieldCategory, "EquipmentStatus")
	ProcLog = NfLog.WithField(logger_util.FieldCategory, "Proc")
	HttpLog = NfLog.WithField(logger_util.FieldCategory, "HTTP")
	UtilLog = NfLog.WithField(logger_util.FieldCategory, "Util")
	SBILog = NfLog.WithField(logger_util.FieldCategory, "SBI")
	DbLog = NfLog.WithField(logger_util.FieldCategory, "DB")
}
