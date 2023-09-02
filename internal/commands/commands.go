package commands

import (
	"context"

	"github.com/Hello-Storage/hello-back/internal/config"
	"github.com/Hello-Storage/hello-back/internal/event"
	"github.com/Hello-Storage/hello-back/internal/server"
)

var log = event.Log

func Start() {
	// init logger
	config.InitLogger()

	// load env
	err := config.LoadEnv()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// connect db and define enum types
	err = config.ConnectDB()
	if err != nil {
		log.Fatal("cannot connect to DB and create enums:", err)
	}

	config.InitDb()

	// connect redis
	config.ConnectRedis()

	// Pass this context down the chain.
	cctx, _ := context.WithCancel(context.Background())

	server.Start(cctx)
}
