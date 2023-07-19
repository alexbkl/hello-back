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
	config, err := config.LoadEnv()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	log.Info("config: file base access key is ", config.FilebaseAccessKey)

	// connect db
	config.ConnectDB()

	// Pass this context down the chain.
	cctx, _ := context.WithCancel(context.Background())

	server.Start(cctx, &config)
}
