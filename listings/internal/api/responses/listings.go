package responses

type ListingResponseBody struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ListingResponse struct {
	Data []ListingResponseBody `json:"data"`
}