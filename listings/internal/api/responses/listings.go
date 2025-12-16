package responses

type ListingResponseBody struct {
	ID             string         `json:"id"`
	Type           string         `json:"type"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Location       string         `json:"location"`
	Endtime        string         `json:"endtime"`
	Seller         SellerBody     `json:"seller"`
	Specifications Specifications `json:"specifications"`
}

type SellerBody struct {
	Name     string `json:"name" `
	Location string `json:"location" `
	ID       string `json:"id" `
}

type Specifications struct {
	Make         string `json:"make"`
	Model        string `json:"model"`
	Engine       string `json:"engine"`
	Colour       string `json:"colour"`
	Transmission string `json:"transmission"`
}

type ListingResponse struct {
	Data []ListingResponseBody `json:"data"`
}
