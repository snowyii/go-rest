package main

import (
	"context"
	"fmt"
	"go-fiber/model"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	app := fiber.New()
	app.Use(recover.New())

	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
			panic(err)
		}
	}()

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
		panic(err)
	}

	fmt.Println("Success connect to mongo")

	db := client.Database("product")

	collection := db.Collection("product")

	app.Get("/product/:id", func(c *fiber.Ctx) error {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		id := c.Params("id")
		var product bson.M
		objID, _ := primitive.ObjectIDFromHex(id)
		fmt.Println((objID))
		if err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{Message: "Product id not found"})
		}
		return c.JSON(product)
	})

	app.Post("/product", func(c *fiber.Ctx) error {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		var product model.Product
		if err := c.BodyParser(&product); err != nil {
			fmt.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{Message: "Bad Request"})

		}
		// _, err := collection.InsertOne(ctx, bson.D{
		// 	{Key: "name", Value: product.Name},
		// 	{Key: "price", Value: product.Price},
		// 	{Key: "description", Value: product.Description},
		// })

		// _, err := collection.InsertOne(ctx, bson.M{
		// 	"name":        product.Name,
		// 	"price":       product.Price,
		// 	"description": product.Description,
		// })
		_, err := collection.InsertOne(ctx, product)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{Message: "Cant create product"})
		}

		return c.JSON(product)
		// var product bson.M

	})

	app.Put("/product/:id", func(c *fiber.Ctx) error {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		id := c.Params("id")
		var product model.Product
		if err := c.BodyParser(&product); err != nil {
			fmt.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{Message: "Bad Request"})

		}
		objID, _ := primitive.ObjectIDFromHex(id)
		fmt.Println(objID)

		update := bson.M{
			"$set": product,
		}
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": objID},
			update,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{Message: "Cant update item"})
		}

		return c.JSON(model.Response{Message: "Update complete"})

	})
	app.Delete("/product/:id", func(c *fiber.Ctx) error {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		id := c.Params("id")
		// var product bson.M
		objID, _ := primitive.ObjectIDFromHex(id)

		res, err := collection.DeleteOne(ctx, bson.M{"_id": objID})

		fmt.Println(err)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{Message: "Error delete item"})
		}
		if res.DeletedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{Message: "Error delete item"})
		}

		return c.JSON(model.Response{Message: "Delete item completed"})
	})

	app.Listen(":8000")
}
