package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/listings/internal/api/requests"
	"github.com/philipjesic/mcg-webapp/listings/internal/api/responses"
	storage "github.com/philipjesic/mcg-webapp/listings/internal/storage/database"
)

type Listings struct {
	db storage.DataStore
}

func CreateListingsHandler(db storage.DataStore) *Listings {
	return &Listings{
		db: db,
	}
}

func (h *Listings) Get(c *gin.Context) {
	ctx := c.Request.Context()
	listings, err := h.db.GetListings(ctx, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Errors: []responses.ErrorMessage{
				{
					Status: http.StatusInternalServerError,
					Title: "internal server error",
					Detail: "failed to fetch listing",
				},
			},
		})
	}
	response := createListingResponse(listings)
	c.JSON(http.StatusOK, response)
}

func (h *Listings) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	listing, err := h.db.GetListingByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Errors: []responses.ErrorMessage{
				{
					Status: http.StatusInternalServerError,
					Title: "internal server error",
					Detail: "failed to fetch listings",
				},
			},
		})
	}

	response := createListingResponse([]storage.Listing{listing})
	c.JSON(http.StatusOK, response)
}

func (h *Listings) Create(c *gin.Context) {
	ctx := c.Request.Context()
	req := requests.ListingCreateRequest{}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	listing := storage.Listing{
		ID:          uuid.New().String(),
		Title:       req.Data.Title,
		Description: req.Data.Description,
	}

	err := h.db.InsertListing(ctx, listing)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Errors: []responses.ErrorMessage{
				{
					Status: http.StatusInternalServerError,
					Title: "internal server error",
					Detail: "failed to create listing",
				},
			},
		})
		return
	}

	res := createListingResponse([]storage.Listing{listing})
	c.JSON(http.StatusCreated, res)
}

func createListingResponse(listings []storage.Listing) responses.ListingResponse {
	listingResponses := make([]responses.ListingResponseBody, 0)
	for _, res := range listings {
		listingResponses = append(listingResponses, createListingResponseBody(res))
	}
	return responses.ListingResponse{
		Data: listingResponses,
	}
}

func createListingResponseBody(listing storage.Listing) responses.ListingResponseBody {
	return responses.ListingResponseBody{
		ID:          listing.ID,
		Type:        "listing",
		Title:       listing.Title,
		Description: listing.Description,
	}
}
