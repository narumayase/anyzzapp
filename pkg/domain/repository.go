package domain

// WhatsAppRepository interface defines the contract for WhatsApp operations
type WhatsAppRepository interface {
	SendMessage(message Message) (*SendMessageResponse, error)
	MarkAsRead(phoneNumberID, messageID string) error
}

// LLMRepository interface defines the contract for the LLM repository responses
type LLMRepository interface {
	SendMessage(prompt string) (string, error)
}
