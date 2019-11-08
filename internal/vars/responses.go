package vars

type ResponseError struct {
	Error string `json:"error"`
}

type Response struct {
	Body interface{} `json:"body"`
}
