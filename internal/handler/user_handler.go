package handler

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"go-api-project/db/sqlc"
	"go-api-project/internal/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type UserHandler struct {
	queries 	*sqlc.Queries
	logger		*zap.Logger
	validate	*validator.Validate
}

func NewUserHandler(
	queries *sqlc.Queries,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		queries: queries,
		logger: logger,
		validate: validator.New(),
	}
}

func calculateAge(dob time.Time) int {
	now := time.Now()

	age := now.Year() - dob.Year()

	if dob.Month() > now.Month() ||
		(dob.Month() == now.Month() &&
			dob.Day() > now.Day()) {
		age--
	}

	return age
}

func (h *UserHandler) GetUsers(c fiber.Ctx) error {
	users, err := h.queries.ListUsers(context.Background())

	if err != nil {
		h.logger.Error(
			"Failed to retrieve users",
			zap.Error(err),
		)

		return c.Status(500).
			SendString("Internal Server Error")
	}

	var response []models.User

	for _, v := range users {
		response = append(response, models.User{
			ID:   v.ID,
			Name: v.Name,
			DOB:  v.Dob.Format("2006-01-02"),
			Age:  calculateAge(v.Dob),
		})
	}

	return c.JSON(response)
}

func (h *UserHandler) GetUser(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(400).
			SendString("Invalid user id")
	}

	user, err := h.queries.GetUser(
		context.Background(),
		int32(id),
	)

	if err == sql.ErrNoRows {
		return c.Status(404).
			SendString("User not found")
	}

	if err != nil {
		h.logger.Error(
			"Failed to retrieve user",
			zap.Int("id", id),
			zap.Error(err),
		)

		return c.Status(500).
			SendString("Internal Server Error")
	}

	return c.JSON(fiber.Map{
		"id":   user.ID,
		"name": user.Name,
		"dob":  user.Dob.Format("2006-01-02"),
		"age":  calculateAge(user.Dob),
	})
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	type PostReq struct {
		Name string `json:"name" validate:"required,min=3,max=50"`
		DOB  string `json:"dob" validate:"required"`
	}

	var req PostReq

	err := c.Bind().Body(&req)

	if err != nil {
		return c.Status(400).
			SendString("Invalid JSON")
	}

	err = h.validate.Struct(req)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"errors": fiber.Map{
				"name": "must be between 3 and 50 characters",
				"dob":  "is required",
			},
		})
	}

	dob, err := time.Parse(
		"2006-01-02",
		req.DOB,
	)

	if err != nil {
		return c.Status(400).
			SendString("Invalid date format")
	}

	res, err := h.queries.CreateUser(
		context.Background(),
		sqlc.CreateUserParams{
			Name: req.Name,
			Dob:  dob,
		},
	)

	if err != nil {
		h.logger.Error(
			"Failed to create user",
			zap.Error(err),
		)

		return c.Status(500).
			SendString("Internal Server Error")
	}

	id, err := res.LastInsertId()

	if err != nil {
		h.logger.Error(
			"Failed to retrieve inserted id",
			zap.Error(err),
		)

		return c.Status(500).
			SendString("Internal Server Error")
	}

	h.logger.Info(
		"User Created",
		zap.Int64("id", id),
	)

	return c.JSON(fiber.Map{
		"id":   id,
		"name": req.Name,
		"dob":  req.DOB,
	})
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(400).
			SendString("Invalid user id")
	}

	type UpdateReq struct {
		Name string `json:"name" validate:"required,min=3,max=50"`
		DOB  string `json:"dob" validate:"required"`
	}

	var req UpdateReq

	err = c.Bind().Body(&req)

	if err != nil {
		return c.Status(400).
			SendString("Invalid JSON")
	}

	err = h.validate.Struct(req)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"errors": fiber.Map{
				"name": "must be between 3 and 50 characters",
				"dob":  "is required",
			},
		})
	}

	dob, err := time.Parse(
		"2006-01-02",
		req.DOB,
	)

	if err != nil {
		return c.Status(400).
			SendString("Invalid date format")
	}

	res, err := h.queries.UpdateUser(
		context.Background(),
		sqlc.UpdateUserParams{
			ID:   int32(id),
			Name: req.Name,
			Dob:  dob,
		},
	)

	if err != nil {
		h.logger.Error(
			"Failed to update user",
			zap.Int("id", id),
			zap.Error(err),
		)

		return c.Status(500).
			SendString("Internal Server Error")
	}

	rows, _ := res.RowsAffected()

	if rows == 0 {
		return c.Status(404).
			SendString("User not found")
	}

	h.logger.Info(
		"User Updated",
		zap.Int("id", id),
	)

	return c.JSON(fiber.Map{
		"id":   id,
		"name": req.Name,
		"dob":  req.DOB,
	})
}

func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(400).
			SendString("Invalid user id")
	}

	res, err := h.queries.DeleteUser(
		context.Background(),
		int32(id),
	)

	if err != nil {
		h.logger.Error(
			"Failed to delete user",
			zap.Int("id", id),
			zap.Error(err),
		)

		return c.Status(500).
			SendString("Internal Server Error")
	}

	rows, _ := res.RowsAffected()

	if rows == 0 {
		return c.Status(404).
			SendString("User not found")
	}

	h.logger.Info(
		"User Deleted",
		zap.Int("id", id),
	)

	return c.SendStatus(204)
}
