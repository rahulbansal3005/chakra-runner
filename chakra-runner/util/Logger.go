package util

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func SetupLogger() {
	// The API for setting attributes is a little different than the package level
	// exported logger. See Godoc.
	Log.Out = os.Stdout

	//You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("chakraRunner.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.Out = file
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}

	Log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}
