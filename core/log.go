package core

import (
	"freemasonry.cc/log"
	"github.com/sirupsen/logrus"
	"strings"
)

//  Lm = LogModel
var (
	LmChainClient       = log.RegisterModule("bc-cli", logrus.DebugLevel)    //
	LmChainType         = log.RegisterModule("bc-ty", logrus.DebugLevel)     // types
	LmChainKeeper       = log.RegisterModule("bc-kp", logrus.DebugLevel)     //keeper
	LmChainCommKeeper   = log.RegisterModule("kp-comm", logrus.DebugLevel)   //comm keeper
	LmChainChatKeeper   = log.RegisterModule("kp-chat", logrus.DebugLevel)   //chat keeper
	LmChainPledgeKeeper = log.RegisterModule("kp-pledge", logrus.DebugLevel) //pledge keeper
	LmChainMsgServer    = log.RegisterModule("bc-ms", logrus.DebugLevel)     //msg server
	LmChainRest         = log.RegisterModule("bc-re", logrus.DebugLevel)     //msg rest
	LmChainMsgAnalysis  = log.RegisterModule("bc-mas", logrus.DebugLevel)    //msg msg analysis
	LmChainUtil         = log.RegisterModule("bc-ut", logrus.DebugLevel)     // util
	LmChainBeginBlock   = log.RegisterModule("bc-bb", logrus.DebugLevel)     // begin block
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
