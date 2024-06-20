
---

# Receipt Service

The Receipt Service is a Go application designed to handle receipt generation and email sending asynchronously using RabbitMQ. It listens for incoming messages on a RabbitMQ queue, generates a PDF receipt based on the provided HTML template, and sends the receipt via email.

## Features

- Asynchronous processing using RabbitMQ
- PDF generation from HTML templates
- Email sending using Gmail SMTP
- Dockerized RabbitMQ setup

## Prerequisites

- Go 1.16+
- Docker and Docker Compose
- RabbitMQ
- Chrome/Chromium for headless PDF generation

## Setup

### RabbitMQ with Docker Compose

To set up RabbitMQ, create a `docker-compose.yml` file with the following content:

```yaml
version: '3.9'
services:
  rabbitmq:
    image: rabbitmq:3-management-alpine
    ports:
      - "5672:5672"
      - "15672:15672"
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    environment:
      RABBITMQ_DEFAULT_USER: "admin"
      RABBITMQ_DEFAULT_PASS: "admin"
volumes:
  rabbitmq_data:
```

Run RabbitMQ using Docker Compose:

```sh
docker-compose up -d
```

### Environment Variables

Create a `.env` file in the root of your project with the following content:

```
PORT=3000
AMQP_URL=amqp://admin:admin@localhost:5672/
GMAIL_EMAIL=your-email@gmail.com
GMAIL_PASSWORD=your-email-password
```

### Project Structure

```plaintext
.
├── api
│   └── handlers
│       └── receipt_handler.go
│   └── routes
│       └── routes.go
├── config
│   └── queue
│       └── rabbitmq_config.go
├── pkg
│   └── email
│       └── email.go
│   └── entities
│       └── email_request.go
├── utils
│   └── pdf.go
├── main.go
├── Dockerfile
├── .env
├── go.mod
└── go.sum
```

### Go Modules

Initialize Go modules and install dependencies:

```sh
go mod init receipt-service
go get github.com/streadway/amqp
go get github.com/gofiber/fiber/v2
go get github.com/chromedp/chromedp
go get github.com/joho/godotenv
```

## Running the Application

To run the application:

```sh
go run main.go
```


## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License.

---

Feel free to modify and expand this README to better fit your project's specific needs.