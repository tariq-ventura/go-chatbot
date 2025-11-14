package logs

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

var Logger *log.Logger

func init() {
	Logger = log.New()

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("No se pudo abrir archivo de logs:", err)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)
	Logger.SetOutput(multiWriter)

	Logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	Logger.SetLevel(log.DebugLevel)
}

func LogInfo(msg string, fields map[string]any) {
	Logger.WithFields(fields).Info(msg)
}

func LogError(msg string, fields map[string]any) {
	Logger.WithFields(fields).Error(msg)
}
