package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"receipt-service/api/routes"
	"receipt-service/config"
	"receipt-service/config/queue"
	"receipt-service/pkg/entities"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Initialize Fiber
	app := fiber.New()

	// Set up routes
	routes.SetupRoutes(app)

	// Start Fiber application
	port := config.GoDotEnvVariable("PORT")
	if port == "" {
		port = "3000" // default port if not specified
	}
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("Error starting server: %s", err)
		}
	}()

	// Set up RabbitMQ connection
	rabbitConfig, err := queue.NewRabbitMQConfig()
	failOnError(err, "Failed to connect to RabbitMQ")
	defer rabbitConfig.Close()

	// Declare queue
	q, err := rabbitConfig.DeclareQueue("email_queue")
	failOnError(err, "Failed to declare a queue")

	// Set up message consumer
	msgs, err := rabbitConfig.Channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var emailReq entities.EmailRequest
			err := json.Unmarshal(d.Body, &emailReq)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}
			fmt.Printf("Email request: %v\n", emailReq)
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	close(forever)
	log.Println("Shutting down...")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
