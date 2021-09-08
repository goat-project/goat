package logger

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

// Init initializes logrus by configuration.
func Init(path string, debug bool) {
	switch path {
	case "":
		if debug {
			InitLogToStdoutDebug()
		} else {
			InitLogToStdout()
		}
	default:
		if debug {
			InitLogToFileDebug(path)
		} else {
			InitLogToFile(path)
		}
	}
}

// InitLogToStdoutDebug inits logrus to log the debug severity or above to Stdout.
func InitLogToStdoutDebug() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

// InitLogToStdout inits logrus to log the info severity or above to Stdout.
func InitLogToStdout() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
}

// InitLogToFile inits logrus to log the info severity or above to the file.
func InitLogToFile(logPath string) {
	logrus.SetFormatter(&logrus.TextFormatter{})

	f, err := os.OpenFile(path.Clean(logPath), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		logrus.Fatalf("error opening file: %v", err)
	}

	logrus.SetOutput(f)
}

// InitLogToFileDebug inits logrus to log the debug severity or above to the file.
func InitLogToFileDebug(logPath string) {
	InitLogToFile(logPath)
	logrus.SetLevel(logrus.DebugLevel)
}
