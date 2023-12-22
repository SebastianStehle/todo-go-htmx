package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"todo/model"

	"todo/views"

	"github.com/a-h/templ"
	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	connStr := "postgresql://postgres:postgres@localhost/todo?sslmode=disable"
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})

	app.Put("/update", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})

	app.Delete("/delete", func(c *fiber.Ctx) error {
		return deleteHandler(c, db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Static("/", "./public")

	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var res string
	var todos []model.Todo
	rows, err := db.Query("SELECT * FROM todos")
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
		c.JSON("An error occured")
	}
	for rows.Next() {
		rows.Scan(&res)
		todo := model.Todo{Item: res}
		todos = append(todos, todo)
	}

	view := views.TodosView(todos)
	return adaptor.HTTPHandler(templ.Handler(view))(c)
}

func postHandler(c *fiber.Ctx, db *sql.DB) error {
	newTodo := model.Todo{}
	if err := c.BodyParser(&newTodo); err != nil {
		log.Printf("An error occured: %v", err)
		return c.SendString(err.Error())
	}
	fmt.Printf("%v", newTodo)
	if newTodo.Item != "" {
		_, err := db.Exec("INSERT into todos VALUES ($1)", newTodo.Item)
		if err != nil {
			log.Fatalf("An error occured while executing query: %v", err)
		}
	}

	view := views.TodoSuccess(newTodo)
	return adaptor.HTTPHandler(templ.Handler(view))(c)
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
	todoToUpdate := c.Query("item")

	updatedTodo := model.Todo{}
	if err := c.BodyParser(&updatedTodo); err != nil {
		log.Printf("An error occured: %v", err)
		return c.SendString(err.Error())
	}

	db.Exec("UPDATE todos SET item=$1 WHERE item=$2", updatedTodo.Item, todoToUpdate)

	view := views.TodoView(updatedTodo)
	return adaptor.HTTPHandler(templ.Handler(view))(c)
}

func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	todoToDelete := c.Query("item")
	db.Exec("DELETE from todos WHERE item=$1", todoToDelete)
	c.Status(http.StatusOK)
	return nil
}
