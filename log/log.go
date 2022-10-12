package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var GTCLog = logrus.New()

func init() {
	GTCLog.SetLevel(logrus.FatalLevel)
	file, err := os.OpenFile("gtc.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		GTCLog.Out = file
	} else {
		GTCLog.Info("Failed to log to file, using default stderr")
	}
}
