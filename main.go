package main

import (
	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/app/auth"
)

func main() {
	inertia, err := NewInertia()
	if err != nil {
		panic(err)
	}

	server := app.NewServer()
	
	auth.New(inertia, server).InitRoutes()

	server.Start(":8080")
}
