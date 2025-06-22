package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/app/auth"
	"github.com/itsbbo/jadel/app/dashboard"
	"github.com/itsbbo/jadel/app/privatekey"
	"github.com/itsbbo/jadel/app/projects"
	"github.com/itsbbo/jadel/app/projects/resources"
	"github.com/itsbbo/jadel/app/repo"
	"github.com/itsbbo/jadel/app/servers"
	"github.com/itsbbo/jadel/app/settings"
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

	database, err := app.NewDB(config)
	if err != nil {
		panic(err)
	}

	defer database.Close()

	server := app.NewServer(config, inertia)

	authRepo := repo.NewAuth(database)
	dashboardRepo := repo.NewDashboard(database)
	settingsRepo := repo.NewSettings(database)
	projectsRepo := repo.NewProject(database)
	serverRepo := repo.NewServer(database)
	privateKeyRepo := repo.NewPrivateKey(database)
	middleware := app.NewMiddleware(server, authRepo, projectsRepo)

	auth.New(server, middleware, authRepo).InitRoutes()
	dashboard.New(server, middleware, dashboardRepo).InitRoutes()
	settings.New(server, middleware, settingsRepo).InitRoutes()
	projects.New(server, middleware, projectsRepo).InitRoutes()
	resources.New(server, middleware, projectsRepo).InitRoutes()
	servers.New(server, middleware, serverRepo).InitRoutes()
	privatekey.New(server, middleware, privateKeyRepo).InitRoutes()

	if config.Server.Debug {
		server.PrintRoutes()
	}

	slog.Info("Server Running", slog.Int("port", config.Server.Port))
	http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), server)
	slog.Info("Server stopped")
}
