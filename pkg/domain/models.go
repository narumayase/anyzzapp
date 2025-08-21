package domain

// Message represents the request to send a message
type Message struct {
	PhoneNumberID string `json:"phone_number_id" binding:"required"`
	To            string `json:"to" binding:"required"`
	Content       string `json:"content" binding:"required"`
	MessageType   string `json:"message_type,omitempty"`
}

// SendMessageResponse represents the entity after sending a message
type SendMessageResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
}

// ErrorResponse represents an error entity
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// WebhookRequest represents incoming webhook data from WhatsApp
type WebhookRequest struct {
	Object string         `json:"object"`
	Entry  []WebhookEntry `json:"entry"`
}

// WebhookEntry represents a single entry in the webhook
type WebhookEntry struct {
	ID      string          `json:"id"`
	Changes []WebhookChange `json:"changes"`
}

// WebhookChange represents a change within a webhook entry
type WebhookChange struct {
	Value WebhookValue `json:"value"`
	Field string       `json:"field"`
}

// WebhookValue represents the value object in webhook changes
type WebhookValue struct {
	MessagingProduct string           `json:"messaging_product"`
	Metadata         WebhookMetadata  `json:"metadata"`
	Contacts         []WebhookContact `json:"contacts,omitempty"`
	Messages         []WebhookMessage `json:"messages,omitempty"`
	Statuses         []WebhookStatus  `json:"statuses,omitempty"`
}

// WebhookMetadata represents metadata in webhook value
type WebhookMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

// WebhookContact represents contact information in webhook
type WebhookContact struct {
	Profile WebhookProfile `json:"profile"`
	WaID    string         `json:"wa_id"`
}

// WebhookProfile represents contact profile information
type WebhookProfile struct {
	Name string `json:"name"`
}

// WebhookMessage represents a message in webhook
type WebhookMessage struct {
	From      string       `json:"from"`
	ID        string       `json:"id"`
	Timestamp string       `json:"timestamp"`
	Text      *WebhookText `json:"text,omitempty"`
	Type      string       `json:"type"`
	// Future message types can be added here
	Image    *WebhookMedia `json:"image,omitempty"`
	Audio    *WebhookMedia `json:"audio,omitempty"`
	Document *WebhookMedia `json:"document,omitempty"`
}

// WebhookText represents text content in a message
type WebhookText struct {
	Body string `json:"body"`
}

// WebhookMedia represents media content in a message
type WebhookMedia struct {
	ID       string `json:"id,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	Caption  string `json:"caption,omitempty"`
	Filename string `json:"filename,omitempty"`
}

// WebhookStatus represents message status updates
type WebhookStatus struct {
	ID           string               `json:"id"`
	Status       string               `json:"status"`
	Timestamp    string               `json:"timestamp"`
	RecipientID  string               `json:"recipient_id"`
	Conversation *WebhookConversation `json:"conversation,omitempty"`
	Pricing      *WebhookPricing      `json:"pricing,omitempty"`
}

// WebhookConversation represents conversation information in status
type WebhookConversation struct {
	ID     string        `json:"id"`
	Origin WebhookOrigin `json:"origin"`
}

// WebhookOrigin represents conversation origin
type WebhookOrigin struct {
	Type string `json:"type"`
}

// WebhookPricing represents pricing information in status
type WebhookPricing struct {
	Billable     bool   `json:"billable"`
	PricingModel string `json:"pricing_model"`
	Category     string `json:"category"`
}
