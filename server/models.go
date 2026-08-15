package server

type WebSocketClientRequest struct {
	Message   string `json:"message"`
	URL       string `json:"url"`
	ErrorCode int    `json:"errorCode"`
}

type RequestPayload struct {
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Headers map[string][]string `json:"headers"`
	Body    []byte              `json:"body"`
}

type ResponsePayload struct {
	StatusCode int                    `json:"status_code"`
	Headers    map[string][]string    `json:"headers"`
	Body       map[string]interface{} `json:"body"`
}
