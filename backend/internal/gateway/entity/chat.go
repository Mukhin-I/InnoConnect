package entity

type WSRequest struct {
    Type   string `json:"type"`
    ChatID int64  `json:"chat_id,omitempty"`
    Text   string `json:"text,omitempty"`
}

type WSResponse struct {
	Type    string `json:"type"`
	ChatID  int64  `json:"chat_id"`
	Message any    `json:"message,omitempty"`
}