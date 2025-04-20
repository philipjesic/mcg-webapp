package requests

type ListingCreateRequestBody struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type ListingCreateRequest struct {
	Data ListingCreateRequestBody `json:"data" binding:"required"`
}
