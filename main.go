package main

import (
	"fmt"
	"log"

	"github.com/saketV8/cine-dots/pkg/database"
	"github.com/saketV8/cine-dots/pkg/handlers"
	"github.com/saketV8/cine-dots/pkg/repositories"
	"github.com/saketV8/cine-dots/pkg/router"

	// log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

// @title           Cine-Dots WatchList API
// @version         1.0
// @description     A watchlist tracker application built with the Gin framework.
// @host            localhost:9090
// @BasePath        /
func main() {
	// log.SetReportCaller(true)

	fmt.Println("============================================")
	fmt.Println("=============== CineDots ===================")
	fmt.Println("============================================")
	fmt.Println()

	// Setting up DB
	db, err := database.InitializeDatabase("sqlite3", "./DB/cine_dots.db")
	if err != nil {
		log.Fatal("DATABASE ERROR: ", err)
	}

	// passing the DB via dependency injection
	app := &router.App{
		WatchListHandler: &handlers.WatchListHandler{
			WatchListModel: &repositories.WatchListModel{
				DB: db.DB,
			},
		},
	}

	router.SetupRouter(app)
}
