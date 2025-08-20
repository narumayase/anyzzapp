# anyzzapp - API de Bot para WhatsApp

Este proyecto proporciona una API que integra la API de Bot de WhatsApp.

## Características

- 📱 **Integración con WhatsApp** - Enviar y recibir mensajes vía WhatsApp Business API
- 🔄 **Soporte de Webhooks** - Manejar mensajes entrantes con una API de LLM

## Requisitos Previos

- Go 1.21 o superior
- Token de acceso a WhatsApp Business API
- Número de teléfono de WhatsApp Business

## Instalación

1. Instalar dependencias:

```bash
go mod tidy
```

2. Configurar variables de entorno:

```bash
cp env.example .env
# Editar .env con los valores descritos abajo.
```

3. Ejecutar la aplicación:

```bash
go run main.go
```

## Configuración

Crear un archivo `.env` basado en `env.example`:

- `WHATSAPP_API_KEY`: Clave de la API de WhatsApp Business (requerido para Meta)
- `WHATSAPP_BASE_URL`: URL de WhatsApp. Ejemplo: https://graph.facebook.com/v18.0
- `WEBHOOK_VERIFY_TOKEN`: Token de verificación del Webhook de WhatsApp.
- `SERVER_PORT`: Puerto del servidor (por defecto: 8080)

### Configuración de WhatsApp Business API

1. **Obtener acceso a la API:**
   - Crear una cuenta de Meta Business
   - Configurar WhatsApp Business API
   - Obtener token de acceso y ID del número de teléfono

2. **Configurar Webhook:**
   - Establecer URL del webhook: `https://tudominio.com/api/v1/whatsapp/webhook`
   - Usar tu `WEBHOOK_VERIFY_TOKEN` para verificación

## 📡 Endpoints

### POST /api/v1/whatsapp/send

**Solicitud:**

```json
{
   "to": "1234567890",
   "content": "¡Hola desde anyzzapp!",
   "message_type": "text"
}
```

**Respuesta:**

```json
{
   "message_id": "1234567890",
   "status": "OK",
   "message": "mensaje de respuesta"
}
```

### Webhook (WhatsApp → Tu API)

```
POST /api/v1/whatsapp/webhook
GET /api/v1/whatsapp/webhook (para verificación)
```

//TODO documentación de webhook

### GET /health

Verifica el estado de la API.

**Respuesta:**

```json
{
   "status": "OK", 
   "message": "anyzzapp API is running"
}
```

#### Usando curl:

```bash
# Verificación de estado
curl http://localhost:8080/health

# Endpoint de chat
curl -X POST http://localhost:8080/api/v1/whatsapp/send \
  -H "Content-Type: application/json" \
  -H "X-Phone-Number-ID: TU_PHONE_NUMBER_ID" \
  -d '{
    "to": "1234567890",
    "content": "¡Hola desde anyzzapp!",
    "message_type": "text"
  }'
```

## 🎗️ Arquitectura

Este proyecto sigue los principios de Clean Architecture:

- **Dominio**: Entidades, interfaces de repositorio y casos de uso
- **Aplicación**: Implementación de casos de uso
- **Infraestructura**: Implementación de repositorios (ej. OpenAI)
- **Interfaces**: Controladores HTTP y routers

## Estructura del Proyecto

```
anyzzapp/
├── cmd/                  # Puntos de entrada de la aplicación
│   └── server/           # Servidor principal
├── internal/             # Código específico del proyecto
│   ├── config/           # Configuración
│   ├── infrastructure/   # Implementaciones de repositorio
│   └── interfaces/       # Controladores HTTP
│       ├── http/         # Handler controller
│       └── middleware/   # Middlewares
├── pkg/                  # Código reutilizable y público
│   ├── domain/           # Entidades e interfaces del dominio
│   └── application/      # Casos de uso
├── main.go               # Punto de entrada principal
├── go.mod                # Dependencias de Go
└── README.md             # Este archivo
```

## Próximos pasos

- **Agregar tu clave de WhatsApp API** en el archivo `.env`
- **Configurar la URL del webhook** en Meta Business Manager

## Backlog

- [ ] Pruebas unitarias
