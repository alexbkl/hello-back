package config

import (
	"sync"

	"github.com/Hello-Storage/hello-back/internal/event"
	"github.com/sirupsen/logrus"
)

var once sync.Once
var log = event.Log

func InitLogger() {
	once.Do(func() {
		log.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
	})
}
