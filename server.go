package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Todo define todo struct
type Todo struct {
	Id        int
	Name      string
	Completed bool
}

var todos = []Todo{
	{Id: 1, Name: "Eat raspberries", Completed: false},
	{Id: 2, Name: "Take out the garbage", Completed: false},
	{Id: 3, Name: "Finish project", Completed: false},
}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World 👋!")
	})

	// Group the endpoint Todo so i don't need to write longer than needed urls
	api := app.Group("/todos")
	api.Get("/", GetTodos)
	api.Post("/", PostTodo)
	api.Get("/:id", GetTodo)
	api.Delete("/:id", DeleteTodo)
	api.Patch("/:id", UpdateTodo)

	app.Get("/todos", GetTodos)

	app.Listen(":5000")
}

func UpdateTodo(ctx *fiber.Ctx) error {
	type request struct {
		Name      *string `json:"name"`
		Completed *bool   `json:"completed"`
	}

	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
	}

	var body request
	err = ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse body",
		})
	}

	var todo *Todo

	for _, value := range todos {
		if value.Id == id {
			todo = &value
			break
		}
	}

	if todo == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Todo not found",
		})
	}

	if body.Name != nil {
		todo.Name = *body.Name
	}

	if body.Completed != nil {
		todo.Completed = *body.Completed
	}

	return ctx.Status(fiber.StatusOK).JSON(todo)
}

func DeleteTodo(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Coudln't parse the id",
		})
	}

	for index, value := range todos {
		if value.Id == id {
			todos = append(todos[0:index], todos[index+1:]...)
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Todo deleted",
			})
		}
	}

	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": "Todo with that Id doesn't exist",
	})
}

func GetTodo(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Couldn't parse the Id",
		})
	}

	for _, value := range todos {
		if value.Id == id {
			return ctx.Status(fiber.StatusOK).JSON(value)
		}
	}

	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": "Todo with that Id doesn't exist",
	})
}

func PostTodo(ctx *fiber.Ctx) error {
	type request struct {
		Name string `json:"name"`
	}

	var body request

	err := ctx.BodyParser(&body)

	// If there was an error while parsing body
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Can't parse body",
		})
	}

	// If the name of the todo is empty
	if body.Name == "" || len(body.Name) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Empty name for a todo",
		})
	}

	// If a todo with that name already exists
	for _, value := range todos {
		if value.Name == body.Name {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Todo with that name already exists",
			})
		}
	}

	todo := Todo{
		Id:        len(todos) + 1,
		Name:      body.Name,
		Completed: false,
	}

	todos = append(todos, todo)

	return ctx.Status(fiber.StatusOK).JSON(todo)
}

func GetTodos(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(todos)
}
