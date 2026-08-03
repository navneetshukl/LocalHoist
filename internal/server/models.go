package server

type WebSocketClientRequest struct {
	Message   string `json:"message"`
	URL       string `json:"url"`
	ErrorCode int    `json:"errorCode"`
}
