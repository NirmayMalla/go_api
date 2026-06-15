package main

import (
	"database/sql"
	"time"

	//"fmt"
	//"go-api-project/internal/models"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v3"
)

func main() {
	
	db, err := sql.Open("mysql", "root:motherjoseph1412*@tcp(localhost:3306)/user_api")
		

	if err != nil {
		log.Fatal(err)		
	}

	// Init a new fiber app
	app := fiber.New()

	// UPDATE A USER
	app.Put("/users/:id", func (c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))

		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid user id\n")
		}
		
		type User struct {
			Name string	`json:"name"` 
			DOB string	`json:"dob"`
		}	
		
		var user User

		err = c.Bind().Body(&user)
		
		if err != nil {
    	return c.Status(400).SendString("Invalid JSON\n")
		}

		res, err := db.Exec("UPDATE users SET name=?, dob=? WHERE id=?", user.Name, user.DOB, id)

		if err != nil {
			return c.Status(500).SendString("Error updating user\n")
		}

		rowsAffected, _ := res.RowsAffected()

		if rowsAffected == 0 {
			return c.Status(400).SendString("User not found\n")
		}
		
		return c.JSON(fiber.Map{
			"id": id,
			"name": user.Name,
			"dob": user.DOB,
		})
	})

	// DELETE A USER
	app.Delete("/users/:id", func (c fiber.Ctx) error {
		
		id, err := strconv.Atoi(c.Params("id"))

		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid user id\n")
		}

		res, err := db.Exec("DELETE FROM users WHERE users.id=?", id)
	
		if err != nil {
			return c.Status(500).SendString("Error deleting user\n")
		}

		rowsAffected, _ := res.RowsAffected()

		if rowsAffected == 0 {
			return c.Status(404).SendString("User not found\n")
		}

		return c.SendStatus(fiber.StatusNoContent)
	
	})

	// ADD A USER
	app.Post("/users", func (c fiber.Ctx) error {

		type User struct {
			Name string	`json:"name"`
			DOB string	`json:"dob"`
		}

		var user User

		err := c.Bind().Body(&user)
		
		if err != nil {
			return c.Status(500).SendString("Error occurred retrieving data\n")
		}
		
		res, err := db.Exec("INSERT INTO users(name, dob) VALUES(?, ?)", user.Name, user.DOB)

		if err != nil {
			return c.Status(400).SendString("Invalid JSON\n")
		}
		
		id, err := res.LastInsertId()

		if err != nil {
			return c.Status(500).SendString("Error occurred retrieving user info\n")
		}
		
		return c.JSON(fiber.Map{
			"id": id,
			"name": user.Name,
			"dob": user.DOB,
		})

	})

	// GET ALL USERS
	app.Get("/users", func (c fiber.Ctx) error {

		current := time.Now()

		rows, err := db.Query("SELECT * FROM users")

		if err != nil {
			return c.Status(500).SendString("Error occurred retrieving users\n")
		}
	
		defer rows.Close()

		var users []fiber.Map

		for rows.Next() {
		
			var id int
			var name string
			var dob string

			err := rows.Scan(&id, &name, &dob)

			if err != nil {
				return c.Status(500).SendString("Error occurred retrieving info\n")
			}
			
			birthDate, err := time.Parse("2006-01-02", dob)	

			if err != nil {
				return c.Status(500).SendString("An error occurred")
			}

			birthYear := birthDate.Year()
			birthMonth := birthDate.Month()
			birthDay := birthDate.Day()

			age := current.Year() - birthYear

			if birthMonth > current.Month() || birthMonth == current.Month() && birthDay > current.Day() {
				age-- 
			}

			users = append(users, fiber.Map{
				"id": id,
				"name": name,
				"dob": dob,
				"age": age,
			})
		}
		return c.JSON(users)
	})

	// GET USER BY ID
	app.Get("/users/:id", func (c fiber.Ctx) error {
		
		current := time.Now()

		id, err := strconv.Atoi(c.Params("id"))
		
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid user id\n");
		}
	
		var userId int
		var name string
		var dob string

		err = db.QueryRow("SELECT * FROM users WHERE users.id=?", id).Scan(&userId, &name, &dob)
					
		if err != nil {
			return c.Status(404).SendString("User not found\n")
		}
		
		birthDate, err := time.Parse("2006-01-02", dob)	

		if err != nil {
			return c.Status(500).SendString("An error occurred")
		}

		birthYear := birthDate.Year()
		birthMonth := birthDate.Month()
		birthDay := birthDate.Day()

		age := current.Year() - birthYear

		if birthMonth > current.Month() || birthMonth == current.Month() && birthDay > current.Day() {
			age-- 
		}

		return c.JSON(fiber.Map{
			"id": userId,
			"name": name,
			"dob": dob,
			"age": age,
		})
	})

	log.Fatal(app.Listen(":3000"))
}


