package requests

type ListingCreateRequestBody struct {
	Title          string         `json:"title" binding:"required"`
	Description    string         `json:"description" binding:"required"`
	Location       string         `json:"location" binding:"required"`
	Endtime        string         `json:"endtime" binding:"required"`
	Seller         SellerBody     `json:"seller" binding:"required"`
	Specifications Specifications `json:"specifications" binding:"required"`
}

type ListingCreateRequest struct {
	Data ListingCreateRequestBody `json:"data" binding:"required"`
}

type SellerBody struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
	ID       string `json:"id" binding:"required"`
}

type Specifications struct {
	Make         string `json:"make" binding:"required"`
	Model        string `json:"model" binding:"required"`
	Engine       string `json:"engine" binding:"required"`
	Colour       string `json:"colour" binding:"required"`
	Transmission string `json:"transmission" binding:"required"`
}
