package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"os"

	"github.com/BernabeSuarez/Apirest-Go/models"
	"github.com/joho/godotenv"
)



func main() {

	app := fiber.New()
	app.Use(cors.New())

	app.Static("/","/public")

	err := godotenv.Load()
  	if err != nil {
    log.Fatal("Error loading .env file")
  }

	mongoUri := os.Getenv("MONGOURI")
	/*Conectarse a mongoDb */
// Use the SetServerAPIOptions() method to set the Stable API version to 1
  serverAPI := options.ServerAPI(options.ServerAPIVersion1)
  opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)
  // Create a new client and connect to the server
  client, err := mongo.Connect(context.TODO(), opts)
  if err != nil {
    panic(err)
  }
   fmt.Println("You successfully connected to MongoDB!",client)


   coll := client.Database("Go-mongo").Collection("tasks")

	
	app.Post("/createtask", func(c *fiber.Ctx) error {
		task := new(models.Task)

		//si hay un error retornar el error
		if err := c.BodyParser(task); err != nil {
            return err
        }
		result, err := coll.InsertOne(context.TODO(), task)
        if err != nil {
   	    panic(err)
        }
		return c.JSON(&fiber.Map{
			"data":result,
		})
	})
	
	app.Get("/todos", func (c *fiber.Ctx) error {
		var tasks  []models.Task
		results, err := coll.Find(context.TODO(), bson.M{})

		if err != nil {
			panic(err)
		}
		//recorremos el resultado y lo guardamos en el array de Tasks
		for results.Next(context.TODO()) {
			var task models.Task
			results.Decode(&task)
			tasks = append(tasks,task)
		}
		c.Status(fiber.StatusOK)
		//devolvemos un array con la informacion de la DB
		return c.JSON(
			tasks,
		)
	
	})

	app.Put("/todos/:id",func(c*fiber.Ctx) error {

		//validar el body
		newTask := new(models.UpdateTask)
		if err := c.BodyParser(newTask); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

		//obtener el id
		id := c.Params("id")
		if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
		}
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
			})
		}
		// actualizar el documento
		result, err := coll.UpdateOne(c.Context(), bson.M{"_id": objectId}, bson.M{"$set": newTask})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update book",
			"message": err.Error(),
		})
	}

	// return the book
	return c.Status(200).JSON(fiber.Map{
		"result": result,
	})
	})

	app.Delete("/todos/:id", func(c*fiber.Ctx) error {
		//obtener el id
		id := c.Params("id")
		if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
		}
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
			})
		}
		//eliminar el documento
		result, err := coll.DeleteOne(c.Context(), bson.M{"_id": objectId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to delete book",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"result": result,
	})
			
	})

	app.Listen(":8080")
	fmt.Println("Server run on port 8080")
}