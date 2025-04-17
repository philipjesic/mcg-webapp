package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	storage "github.com/philipjesic/listings/storage/database"
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

func CreateListingsHandler(db storage.DataStore) *Listings {
	return &Listings{
		db: db,
	}
}

func (h *Listings) Get(c *gin.Context) {
	id := c.Param("id")
	listing, err := h.db.GetListingByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error...")
	}

	c.JSON(http.StatusOK, listing)
}

func (h *Listings) Create(c *gin.Context) {
	req := ListingCreateRequest{}

	c.ShouldBindBodyWithJSON(&req)

	listing := storage.Listing{
		Title:       req.Data.Title,
		Description: req.Data.Description,
	}

	err := h.db.InsertListing(listing)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error...")
		return
	}

	c.JSON(http.StatusCreated, listing)
}
