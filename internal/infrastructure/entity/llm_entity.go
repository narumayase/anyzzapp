package entity

type Request struct {
	Prompt string `json:"prompt"`
}

type Response struct {
	Response string `json:"response"`
}
