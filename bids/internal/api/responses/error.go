package responses

type ErrorResponse struct {
	Errors []ErrorMessage `json:"errors"`
}

type ErrorMessage struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}
