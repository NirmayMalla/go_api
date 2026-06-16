package main

import (
	"log"

	"go-api-project/config"
	"go-api-project/db/sqlc"
	"go-api-project/internal/logger"
	"go-api-project/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func main() {

	appLogger, err := logger.NewLogger()

	if err != nil {
		log.Fatal(err)
	}

	defer appLogger.Sync()
	
	db, err := config.NewDB()
		
	if err != nil {
		log.Fatal(err)		
	}

	defer db.Close()

	queries := sqlc.New(db)

	// Init a new fiber app
	app := fiber.New()

	routes.Setup(
		app,
		queries,
		appLogger,
	)

	log.Fatal(
		app.Listen(":3000"),
	)
}


