package controllers

import (
	"github.com/MishraShardendu22/schema"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func MakeBlog(c *fiber.Ctx, collections *mongo.Collection) error {
	var blog schema.Post

	if err := c.BodyParser(&blog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error parsing blog data"})
	}

	_, err := collections.InsertOne(c.Context(), blog)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error adding blog"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Blog added successfully"})
}

func DeleteBlog(c *fiber.Ctx, collections *mongo.Collection) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id is required"})
	}

	_, err := collections.DeleteOne(c.Context(), bson.M{"_id": id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error deleting blog"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Blog deleted successfully"})
}

func EditBlog(c *fiber.Ctx, collections *mongo.Collection) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id is required"})
	}

	var blog schema.Post
	if err := c.BodyParser(&blog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error parsing blog data"})
	}

	update := bson.M{"$set": blog}
	_, err := collections.UpdateOne(c.Context(), bson.M{"_id": id}, update)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error updating blog"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Blog updated successfully"})
}

func GetBlog(c *fiber.Ctx, collections *mongo.Collection) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id is required"})
	}

	cursor, err := collections.Find(c.Context(), bson.M{"_id": id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error finding blog"})
	}
	defer cursor.Close(c.Context())

	var blogs []schema.Post
	if err := cursor.All(c.Context(), &blogs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error decoding blog data"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Blog retrieved successfully", "data": blogs})
}