package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

/*
ConfigureLogger sets up the logger with a specific format and output.
We can configure the logger as per our requirements.
*/
func ConfigureLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)
	return logger
}
