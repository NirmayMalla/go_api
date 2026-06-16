package routes

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"go-api-project/db/sqlc"
	"go-api-project/internal/handler"
)

func Setup(
	app *fiber.App,
	queries *sqlc.Queries,
	logger *zap.Logger,
)	{
		userHandler := handler.NewUserHandler(
			queries,
			logger,
		)

		app.Get(
			"/users",
			userHandler.GetUsers,
		)

		app.Get(
			"/users/:id",
			userHandler.GetUser,
		)

		app.Post(
			"/users",
			userHandler.CreateUser,
		)	

		app.Put(
			"/users/:id",
			userHandler.UpdateUser,
		)

		app.Delete(
			"/users/:id",
			userHandler.DeleteUser,
		)
}
