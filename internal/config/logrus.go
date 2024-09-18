package config

import (
	"fmt"
	"github.com/bccfilkom/career-path-service/internal/pkg/env"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

func NewLogger() *logrus.Logger {
	fileName := fmt.Sprintf("logs/start-on-%d-log.log", time.Now().UnixMilli())

	f, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0777)
	log := logrus.New()
	if env.GetBool("PRODUCTION", false) {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(logrus.DebugLevel)
	}
	log.SetFormatter(&logrus.JSONFormatter{})
	log.Out = io.MultiWriter(os.Stderr, f)
	log.SetReportCaller(true)

	return log
}
