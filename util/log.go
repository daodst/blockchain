package util

import (
	"freemasonry.cc/log"
	"github.com/sirupsen/logrus"
	"strings"
)

//  Lm = LogModel
var (
	LmChainClient      = log.RegisterModule("ccli", logrus.InfoLevel)      //
	LmChainType        = log.RegisterModule("chain-t", logrus.InfoLevel)   // types
	LmChainKeeper      = log.RegisterModule("chain-kp", logrus.InfoLevel)  //keeper
	LmChainMsgServer   = log.RegisterModule("chain-ms", logrus.InfoLevel)  //msg server
	LmChainRest        = log.RegisterModule("chain-re", logrus.InfoLevel)  //msg rest
	LmChainMsgAnalysis = log.RegisterModule("chain-mas", logrus.InfoLevel) //msg msg analysis
	LmChainUtil        = log.RegisterModule("chain-ut", logrus.InfoLevel)  // util
)

//
func BuildLog(funcName string, modules ...log.LogModule) *logrus.Entry {
	moduleName := ""
	for _, v := range modules {
		if moduleName != "" {
			moduleName += "/"
		}
		moduleName += string(v)
	}
	logEntry := log.Log.WithField("module", strings.ToLower(moduleName))
	if funcName != "" {
		logEntry = logEntry.WithField("method", strings.ToLower(funcName))
	}
	return logEntry
}
