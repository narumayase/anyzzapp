# anyzzapp - WhatsApp Bot API

This project provides an API that integrates WhatsApp Bot API and responds to the message with an LLM API.

## Features

- 📱 **WhatsApp Integration** - Send and receive messages via WhatsApp Business API.
- 🔄 **Webhook Support** - Handle incoming messages with an LLM API.

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

- `WHATSAPP_API_KEY`: WhatsApp Business API key (required for Meta, view https://developers.facebook.com)
- `WHATSAPP_BASE_URL`: WhatsApp URL. Example: https://graph.facebook.com/v20.0
- `WEBHOOK_VERIFY_TOKEN`: WhatsApp Webhook Verify.
- `SERVER_PORT`: Server port (default: 8080)
- `LLM_URL`: LLM API URL

### WhatsApp Business API Setup

https://developers.facebook.com

1. **Get API Access:**
    - Create a Meta Business account
    - Set up WhatsApp Business API
    - Get your access token and phone number ID

2. **Configure Webhook:**
    - Set webhook URL in Meta: `https://yourdomain.com/api/v1/whatsapp/webhook`
    - Use your `WEBHOOK_VERIFY_TOKEN` for verification

## 📡 Endpoints

### POST /api/v1/whatsapp/send

**Request:**

```json
{
   "to": "541112345678", 
   "content": "Hello from anyzzapp!", 
   "message_type": "text"
}
```

to: is the recipient's phone number in E.164 international format (without +, without spaces, without dashes).

**Response:**

```json
{
   "message_id": "1234567890", 
   "status": "OK",
   "message": "response message"
}
```

### Webhook (WhatsApp → Your API)

https://developers.facebook.com/docs/whatsapp/cloud-api/get-started#configure-webhooks

```
POST /api/v1/whatsapp/webhook
```

### Verification endpoint (used by WhatsApp Business API):

```
GET /api/v1/whatsapp/webhook
```

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
  -H "X-Phone-Number-ID: YOUR_BOT_PHONE_NUMBER_ID" \
  -d '{
    "to": "541112345678",
    "content": "Hello from anyzzapp!",
    "message_type": "text"
  }'
```

note: in test mode, WhatsApp does not allow outgoing text messages from your bot, first the user must write to the bot so that it can reply.

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
├── README_ES.md          # README in spanish
└── README.md             # This file
```

## Next Steps

- **Add your WhatsApp API key** to the `.env` file
- **Configure webhook URL** in Meta Business Manager

## Backlog

- [ ] Unit Tests 
