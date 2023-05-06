package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"todo-app/ent"
	"todo-app/graph"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func graphqlHandler(client *ent.Client) fiber.Handler {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Client: client}}))
	httpHandler := fasthttpadaptor.NewFastHTTPHandler(h)

	return func(c *fiber.Ctx) error {
		httpHandler(c.Context())
		return nil
	}
}

// Open new connection
func Open(databaseUrl string) *ent.Client {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}

func main() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Create the PostgreSQL connection string
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	client := Open(connString)

	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	app := fiber.New()
	app.Get("/todos", getTodos(client))
	app.Post("/todos", createTodo(client))
	app.Get("/todos/:id", getTodo(client))
	app.Put("/todos/:id", updateTodo(client))
	app.Delete("/todos/:id", deleteTodo(client))
	app.Post("/graphql", graphqlHandler(client))

	app.Listen(":3000")
}

func getTodos(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		todos, err := client.Todo.Query().All(context.Background())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot retrieve todos",
			})
		}
		return c.JSON(todos)
	}
}

func createTodo(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := new(ent.Todo)
		if err := c.BodyParser(t); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}
		created, err := client.Todo.Create().SetTask(t.Task).SetCompleted(t.Completed).Save(context.Background())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot create todo",
			})
		}
		return c.Status(fiber.StatusCreated).JSON(created)
	}
}

func getTodo(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}
		t, err := client.Todo.Get(context.Background(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Todo not found",
			})
		}
		return c.JSON(t)
	}
}

func updateTodo(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}
		t := new(ent.Todo)
		if err := c.BodyParser(t); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		updater := client.Todo.UpdateOneID(id).SetTask(t.Task).SetCompleted(t.Completed)
		updated, err := updater.Save(context.Background())
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Todo not found",
			})
		}
		return c.JSON(updated)
	}
}

func deleteTodo(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}
		err = client.Todo.DeleteOneID(id).Exec(context.Background())
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Todo not found",
			})
		}
		return c.JSON(fiber.Map{
			"result": "Todo deleted",
		})
	}
}
