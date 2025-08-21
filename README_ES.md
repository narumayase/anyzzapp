# anyzzapp - API de Bot para WhatsApp + LLM

Este proyecto proporciona una API que integra la API de Bots de WhatsApp y responde el mensaje con una LLM API.

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

- `WHATSAPP_API_KEY`: Clave de la API de WhatsApp Business (requerido para Meta, ver https://developers.facebook.com)
- `WHATSAPP_BASE_URL`: URL de WhatsApp. Ejemplo: https://graph.facebook.com/v20.0
- `WEBHOOK_VERIFY_TOKEN`: Token de verificación del Webhook de WhatsApp.
- `SERVER_PORT`: Puerto del servidor (por defecto: 8080)
- `LLM_URL`: LLM API URL

### Configuración de WhatsApp Business API

1. **Obtener acceso a la API:**
   - Crear una cuenta de Meta Business
   - Configurar WhatsApp Business API
   - Obtener token de acceso y ID del número de teléfono

2. **Configurar Webhook:**
   - Establecer URL del webhook en Meta: `https://tudominio.com/api/v1/whatsapp/webhook`
   - Usar tu `WEBHOOK_VERIFY_TOKEN` para verificación

## 📡 Endpoints

### POST /api/v1/whatsapp/send

**Solicitud:**

```json
{
   "to": "541112345678",
   "content": "¡Hola desde anyzzapp!",
   "message_type": "text"
}
```

to: es el número de teléfono del destinatario en formato internacional E.164 (sin +, sin espacios, sin guiones).

**Respuesta:**

```json
{
   "message_id": "1234567890",
   "status": "OK",
   "message": "mensaje de respuesta"
}
```

### Webhook (WhatsApp → Esta API)

https://developers.facebook.com/docs/whatsapp/cloud-api/get-started#configure-webhooks

```
POST /api/v1/whatsapp/webhook
```

#### Endpoint de Verificación del Webhook (usado por WhatsApp Business API):

```
GET /api/v1/whatsapp/webhook
```

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
    "to": "541112345678",
    "content": "¡Hola desde anyzzapp!",
    "message_type": "text"
  }'
```

nota: en modo prueba, WhatsApp no permite mensajes salientes del tipo `text` de su bot, primero el usuario debe escribirle al bot para que este pueda responder.

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
├── README_ES.md          # Este archivo
└── README.md             # README en inglés
```

## 🧪 Pruebas

### Ejecutar Pruebas

Para ejecutar todas las pruebas:

```bash
go test ./...
```

Para ejecutar pruebas con salida detallada:

```bash
go test -v ./...
```

Para ejecutar pruebas de un paquete específico:

```bash
go test ./internal/config/
go test ./cmd/server/
```

### Cobertura de Pruebas

Para verificar la cobertura de pruebas (excluyendo mocks):

```bash
# Generar reporte de cobertura
go test -coverprofile=coverage.out ./...

# Ver reporte de cobertura en terminal
go tool cover -func=coverage.out

# Generar reporte HTML de cobertura
go tool cover -html=coverage.out -o coverage.html

# Ver cobertura excluyendo mocks
go test -coverprofile=coverage.out ./... && \
go tool cover -func=coverage.out | grep -v "mocks"
```

### Ejecutar Benchmarks

```bash
go test -bench=. ./...
```

## Próximos pasos

- **Agregar tu clave de WhatsApp API** en el archivo `.env`
- **Configurar la URL del webhook** en Meta Business Manager

## Backlog

- [x] Pruebas unitarias
- [ ] Pruebas de integración
- [ ] Documentación de API con Swagger