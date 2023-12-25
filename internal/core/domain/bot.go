package domain

type TextMessageRequest struct {
	Message string `json:"message"`
}

type TextMessageResponse struct {
	ReplyMessage string `json:"message"`
}
