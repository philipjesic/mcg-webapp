package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	storage "github.com/philipjesic/mcg-webapp/listings/storage/database"
)

type Listings struct {
	db storage.DataStore
}

type ListingCreateRequestBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ListingCreateRequest struct {
	Data ListingCreateRequestBody `json:"data"`
}

type ListingResponseBody struct {
	ID          string `json:"id"`
	Type        string ` json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ListingResponse struct {
	Data []ListingResponseBody `json:"data"`
}

func CreateListingsHandler(db storage.DataStore) *Listings {
	return &Listings{
		db: db,
	}
}

func (h *Listings) Get(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	listing, err := h.db.GetListingByID(ctx, id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error...")
	}

	response := createListingResponse([]storage.Listing{listing})
	c.JSON(http.StatusOK, response)
}

func (h *Listings) Create(c *gin.Context) {
	ctx := c.Request.Context()
	req := ListingCreateRequest{}

	c.ShouldBindBodyWithJSON(&req)

	listing := storage.Listing{
		ID:          uuid.New().String(),
		Title:       req.Data.Title,
		Description: req.Data.Description,
	}

	err := h.db.InsertListing(ctx, listing)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Internal Server Error...")
		return
	}

	res := createListingResponse([]storage.Listing{listing})
	c.JSON(http.StatusCreated, res)
}

func createListingResponse(listings []storage.Listing) ListingResponse {
	listingResponses := make([]ListingResponseBody, 0)
	for _, res := range listings {
		listingResponses = append(listingResponses, createListingResponseBody(res))
	}
	return ListingResponse{
		Data: listingResponses,
	}
}

func createListingResponseBody(listing storage.Listing) ListingResponseBody {
	return ListingResponseBody{
		ID:          listing.ID,
		Type:        "listing",
		Title:       listing.Title,
		Description: listing.Description,
	}
}
