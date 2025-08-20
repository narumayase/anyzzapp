# anyzzapp - WhatsApp Bot API

This project provides an API that integrates WhatsApp Bot API.

## Features

- 📱 **WhatsApp Integration** - Send and receive messages via WhatsApp Business API
- 🔄 **Webhook Support** - Handle incoming messages with an LLM API

## Prerequisites

- Go 1.21 or higher
- WhatsApp Business API access token
- WhatsApp Business phone number

## Installation

1. Install dependencies:

```bash
go mod tidy
```

2. Configure environment variables:

```bash
cp env.example .env
# Edit .env with the values described below.
```

3. Run the application:

```bash
go run main.go
```

## Configuration

Create a `.env` file based on `env.example`:

- `WHATSAPP_API_KEY`: WhatsApp Business API key (required for Meta)
- `WHATSAPP_BASE_URL`: WhatsApp URL. Example: https://graph.facebook.com/v18.0
- `WEBHOOK_VERIFY_TOKEN`: WhatsApp Webhook Verify. 
- `SERVER_PORT`: Server port (default: 8080)

### WhatsApp Business API Setup

1. **Get API Access:**
    - Create a Meta Business account
    - Set up WhatsApp Business API
    - Get your access token and phone number ID

2. **Configure Webhook:**
    - Set webhook URL: `https://yourdomain.com/api/v1/whatsapp/webhook`
    - Use your `WEBHOOK_VERIFY_TOKEN` for verification

## 📡 Endpoints

### POST /api/v1/whatsapp/send

**Request:**

```json
{
   "to": "1234567890", 
   "content": "Hello from anyzzapp!", 
   "message_type": "text"
}
```

**Response:**

```json
{
   "message_id": "1234567890", 
   "status": "OK",
   "message": "response message"
}
```

### Webhook (WhatsApp → Your API)

```
POST /api/v1/whatsapp/webhook
GET /api/v1/whatsapp/webhook (for verification)
```

//TODO webhook docs

### GET /health

Checks the API status.

**Response:**

```json
{
  "status": "OK",
  "message": "anyzzapp API is running"
}
```

#### Using curl:

```bash
# Health check
curl http://localhost:8080/health

# Chat endpoint
curl -X POST http://localhost:8080/api/v1/whatsapp/send \
  -H "Content-Type: application/json" \
  -H "X-Phone-Number-ID: YOUR_PHONE_NUMBER_ID" \
  -d '{
    "to": "1234567890",
    "content": "Hello from anyzzapp!",
    "message_type": "text"
  }'
```

## 🎗️ Architecture

This project follows Clean Architecture principles:

- **Domain**: Entities, repository interfaces, and use cases
- **Application**: Implementation of use cases
- **Infrastructure**: OpenAI repository implementation
- **Interfaces**: HTTP controllers and routers

## Project Structure

```
anyzzapp/
├── cmd/                  # Application entry points
│   └── server/           # Main server
├── internal/             # Project-specific code
│   ├── config/           # Configuration
│   ├── infrastructure/   # Repository implementations
│   └── interfaces/       # HTTP controllers
│       ├── http/         # Handler controller
│       └── middleware/   # Middlewares
├── pkg/                  # Reusable and public code
│   ├── domain/           # Domain entities and interfaces
│   └── application/      # Use cases
├── main.go               # Main entry point
├── go.mod                # Go dependencies
└── README.md             # This file
```

## Next Steps

- **Add your WhatsApp API key** to the `.env` file
- **Configure webhook URL** in Meta Business Manager

## Backlog

- [ ] Unit Tests 
