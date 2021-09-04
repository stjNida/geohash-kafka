package main

import (
	"encoding/json"
	"kafka-golang/producer"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()
	api := app.Group("/api/v1")

	api.Post("/geolocation", createGeolocation)

	app.Listen(":3000")

}

func createGeolocation(c *fiber.Ctx) error {

	// Instantiate new Message struct
	cmt := new(producer.Geolocation)

	//  Parse body into comment struct
	if err := c.BodyParser(cmt); err != nil {
		log.Println(err)
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}
	// convert body into bytes and send it to kafka
	cmtInBytes, err := json.Marshal(cmt)
	producer.PushToQueue("geolocations", cmtInBytes)

	// Return Comment in JSON format
	err = c.JSON(&fiber.Map{
		"success": true,
		"message": "Comment pushed successfully",
		"comment": cmt,
	})
	if err != nil {
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Error creating product",
		})
		return err
	}

	return err
}
