package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

/*
ConfigureLogger sets up the logger with a specific format and output.
We can configure the logger as per our requirements.
*/
func ConfigureLogger() *logrus.Logger {
	if Log == nil {
		Log = logrus.New()
		Log.SetLevel(logrus.InfoLevel)
		Log.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
		Log.SetOutput(os.Stdout)
	}
	return Log
}
