package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"receipt-service/api/handlers"
	"receipt-service/api/routes"
	"receipt-service/config"
	"receipt-service/config/queue"
	"receipt-service/pkg/email"
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
		false, // auto-ack to false,
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
				d.Ack(false)
				continue
			}
			fmt.Printf("Received email request: %+v\n", emailReq)

			datum, err := json.Marshal(emailReq.Data)
			if err != nil {
				log.Fatalf("Error marshaling data: %v", err)
			}

			// Print email request
			outputPath, err := handlers.PrintReceipt(emailReq.ReceiptType, datum)
			if err != nil {
				log.Printf("Failed to print receipt: %v", err)
				d.Nack(false, true) // Requeue the message for retry
				continue
			}
			// Send email
			smtpHost := "smtp.gmail.com"
			smtpPort := "587"
			senderEmail := config.GoDotEnvVariable("GMAIL_EMAIL")
			senderPassword := config.GoDotEnvVariable("GMAIL_PASSWORD")
			fmt.Printf("email := %s\n", senderEmail)
			emailService := email.NewSMTPEmailService(smtpHost, smtpPort, senderEmail, senderPassword)

			err = emailService.SendEmail(emailReq.UserEmail, outputPath)
			if err != nil {
				log.Fatalf("Failed to send email: %v", err)
				d.Nack(false, true) // Requeue the message for retry
			}
			// remove the message from the queue
			d.Ack(false)
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
