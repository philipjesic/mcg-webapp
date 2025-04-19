package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	storage "github.com/philipjesic/mcg-webapp/listings/storage/database"
)

func Test_GetListing_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := new(storage.MockDataStore)
	handler := CreateListingsHandler(mockDB)

	listing := storage.Listing{
		ID:          "listing-id-123",
		Title:       "Sample Listing",
		Description: "A test listing for unit test",
	}

	mockDB.On("GetListingByID", mock.Anything, "listing-id-123").Return(listing, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/listings/listing-id-123", nil)
	rec := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/listings/:id", handler.Get)

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockDB.AssertExpectations(t)
}

func Test_GetListing_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := new(storage.MockDataStore)
	handler := CreateListingsHandler(mockDB)

	// Simulate DB error
	mockDB.On("GetListingByID", mock.Anything, "fail-id").Return(storage.Listing{}, errors.New("db failed")).Once()

	req, _ := http.NewRequest(http.MethodGet, "/listings/fail-id", nil)
	rec := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/listings/:id", handler.Get)

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockDB.AssertExpectations(t)
}

func Test_CreateListing_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := new(storage.MockDataStore)
	handler := CreateListingsHandler(mockDB)

	requestBody := `{
		"data": {
			"title": "Test Title",
			"description": "Test Description"
		}
	}`

	// This is tricky: we want to capture the listing passed into InsertListing
	var insertedListing storage.Listing
	mockDB.On("InsertListing", mock.Anything, mock.MatchedBy(func(l storage.Listing) bool {
		insertedListing = l
		return l.Title == "Test Title" && l.Description == "Test Description"
	})).Return(nil).Once()

	req, _ := http.NewRequest(http.MethodPost, "/listings", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r := gin.Default()
	r.POST("/listings", handler.Create)

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	// Parse response JSON
	var res ListingResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Len(t, res.Data, 1)
	assert.Equal(t, insertedListing.ID, res.Data[0].ID)
	assert.Equal(t, "listing", res.Data[0].Type)
	assert.Equal(t, insertedListing.Title, res.Data[0].Title)
	assert.Equal(t, insertedListing.Description, res.Data[0].Description)

	// ID should be valid UUID
	_, err = uuid.Parse(res.Data[0].ID)
	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

func Test_CreateListing_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := new(storage.MockDataStore)
	handler := CreateListingsHandler(mockDB)

	requestBody := `{
		"data": {
			"title": "Fail Insert",
			"description": "This should trigger a DB error"
		}
	}`

	mockDB.On("InsertListing", mock.Anything, mock.AnythingOfType("storage.Listing")).
		Return(errors.New("simulated insert error")).Once()

	req, _ := http.NewRequest(http.MethodPost, "/listings", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r := gin.Default()
	r.POST("/listings", handler.Create)

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "Internal Server Error...", rec.Body.String())

	mockDB.AssertExpectations(t)
}