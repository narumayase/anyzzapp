package entity

// SendMessagePayload represents the payload structure for sending messages
type SendMessagePayload struct {
	MessagingProduct string `json:"messaging_product"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             *Text  `json:"text,omitempty"`
}

type Text struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}

// Result message response
type Result struct {
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}
