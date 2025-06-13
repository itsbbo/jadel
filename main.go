package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/app/auth"
	"github.com/itsbbo/jadel/app/dashboard"
	"github.com/itsbbo/jadel/app/repo"
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

	authRepo := repo.NewAuth(database)
	dashboardRepo := repo.NewDashboard(database)
	server := app.NewServer(config, inertia)

	auth.New(server, authRepo).InitRoutes()
	dashboard.New(server, dashboardRepo).InitRoutes()

	if config.Server.Debug {
		server.PrintRoutes()
	}

	slog.Info("Server Running", slog.Int("port", config.Server.Port))
	http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), server)
	slog.Info("Server stopped")
}
