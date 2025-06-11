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

	database, err := app.NewPostgres(config)
	if err != nil {
		panic(err)
	}

	defer database.Close()

	server := app.NewServer()

	auth.New(inertia, server, database).InitRoutes()

	server.Start(fmt.Sprintf(":%d", config.Server.Port))
}
