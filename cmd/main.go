package main

import (
	"context"
	"fmt"
	"go-fiber/model"
	"log"
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

	uri := "mongodb+srv://admin:poiuytrewq@oxbidkrub.zgesmxt.mongodb.net/?retryWrites=true&w=majority"

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
			return c.JSON("No document")
		}
		return c.JSON(product)
	})

	app.Post("/product", func(c *fiber.Ctx) error {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		var product model.Product
		if err := c.BodyParser(&product); err != nil {
			fmt.Print(err)
			return c.JSON("Failed to bind")

		}
		_, err := collection.InsertOne(ctx, bson.D{
			{Key: "name", Value: product.Name},
			{Key: "price", Value: product.Price},
			{Key: "description", Value: product.Description},
		})

		if err != nil {
			return c.JSON("Error write to DB")
		}

		return c.JSON("Write Complete")
		// var product bson.M

	})

	app.Put("/product/:id", func(c *fiber.Ctx) error {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		id := c.Params("id")
		var product model.Product
		if err := c.BodyParser(&product); err != nil {
			fmt.Print(err)
			return c.JSON("Failed to bind")

		}
		objID, _ := primitive.ObjectIDFromHex(id)
		fmt.Println(objID)

		update := bson.D{
			{"$set", bson.D{{"name", product.Name}}},
			{"$set", bson.D{{"price", product.Price}}},
			{"$set", bson.D{{"description", product.Description}}},
		}
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": objID},
			update,
		)
		if err != nil {
			return c.JSON("Cant update")
		}

		return c.JSON("Update complete")

	})
	app.Delete("/product/:id", func(c *fiber.Ctx) error {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		id := c.Params("id")
		// var product bson.M
		objID, _ := primitive.ObjectIDFromHex(id)

		fmt.Println(objID)
		if _, err := collection.DeleteOne(ctx, bson.M{"_id": objID}); err != nil {
			return c.JSON("No document")
		}

		return c.JSON("Finish delete")
	})

	app.Listen(":8000")
}
