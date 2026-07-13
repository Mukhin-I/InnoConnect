package entity

// WSRequest represents a WebSocket request message
type WSRequest struct {
	Type   string `json:"type"`
	ChatID int64  `json:"chat_id,omitempty"`
	Text   string `json:"text,omitempty"`
}

// WSResponse represents a WebSocket response message
type WSResponse struct {
	Type    string `json:"type"`
	ChatID  int64  `json:"chat_id"`
	Message any    `json:"message,omitempty"`
}
