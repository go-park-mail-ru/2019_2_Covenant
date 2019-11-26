package vars

type Body map[string]interface{}

type Response struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Body    *Body  `json:"body,omitempty"`
}
