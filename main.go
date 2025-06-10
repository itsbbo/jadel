package main

import (
	"fmt"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/app/auth"
)

func main() {
	config, err := app.InitConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	inertia, err := app.NewInertia()
	if err != nil {
		panic(err)
	}

	server := app.NewServer()

	auth.New(inertia, server).InitRoutes()

	server.Start(fmt.Sprintf(":%d", config.Server.Port))
}
