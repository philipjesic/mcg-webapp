package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/philipjesic/mcg-webapp/bids/internal/api/responses"
	"github.com/philipjesic/mcg-webapp/bids/internal/api/routes"
	"github.com/philipjesic/mcg-webapp/bids/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetBid_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockDB := new(storage.MockDataStore)
	routes.RegisterAPI(r, mockDB)

	bid := storage.Bid{
		AuctionID: "test-auction-id",
		ID:        "test-bid-id",
		UserID:    "test-user-id",
		Amount:    10000,
		Timestamp: time.Now(),
	}

	mockDB.On("GetBidByID", mock.Anything, "test-bid-id").Return(bid, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/bids/test-bid-id", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockDB.AssertExpectations(t)
}

func Test_GetBid_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockDB := new(storage.MockDataStore)
	routes.RegisterAPI(r, mockDB)

	// Simulate DB error
	mockDB.On("GetBidByID", mock.Anything, "fail-id").Return(storage.Bid{}, errors.New("db failed")).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/bids/fail-id", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockDB.AssertExpectations(t)
}

func Test_CreateListing_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockDB := new(storage.MockDataStore)
	routes.RegisterAPI(r, mockDB)
	timeStamp := time.Now().UTC()

	requestBody := `{
		"data": {
			"auction_id": "test-auction",
			"user_id": "test-user",
			"amount": 7000,
			"timestamp": "` + timeStamp.Format(time.RFC3339Nano) + `"
		}
	}`

	var insertedBid storage.Bid
	mockDB.On("InsertBid", mock.Anything, mock.MatchedBy(func(b storage.Bid) bool {
		insertedBid = b
		return b.AuctionID == "test-auction" && b.UserID == "test-user" &&
			b.Amount == 7000 && b.Timestamp == timeStamp
	})).Return(nil).Once()

	req, _ := http.NewRequest(http.MethodPost, "/api/bids", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	// Parse response JSON
	var res responses.BidsResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Len(t, res.Data, 1)
	assert.Equal(t, insertedBid.ID, res.Data[0].ID)
	assert.Equal(t, "test-auction", res.Data[0].AuctionID)
	assert.Equal(t, "test-user", res.Data[0].UserID)
	assert.Equal(t, 7000, res.Data[0].Amount)
	assert.Equal(t, timeStamp.String(), res.Data[0].Timestamp.String())

	// ID should be valid UUID
	_, err = uuid.Parse(res.Data[0].ID)
	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

func Test_CreateListing_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockDB := new(storage.MockDataStore)
	routes.RegisterAPI(r, mockDB)
	timeStamp := time.Now().UTC()

	requestBody := `{
		"data": {
			"auction_id": "test-auction",
			"user_id": "test-user",
			"amount": 7000,
			"timestamp": "` + timeStamp.Format(time.RFC3339Nano) + `"
		}
	}`

	mockDB.On("InsertBid", mock.Anything, mock.MatchedBy(func(b storage.Bid) bool {
		return b.AuctionID == "test-auction" && b.UserID == "test-user" &&
			b.Amount == 7000 && b.Timestamp == timeStamp
	})).Return(errors.New("simulated insert error")).Once()

	req, _ := http.NewRequest(http.MethodPost, "/api/bids", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "{\"errors\":[{\"status\":500,\"title\":\"internal server error\",\"detail\":\"failed to create bid\"}]}", rec.Body.String())

	mockDB.AssertExpectations(t)
}

func Test_CreateListing_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockDB := new(storage.MockDataStore)
	routes.RegisterAPI(r, mockDB)

	body := []byte(`{
		"data": {
		}
	}`)

	req, _ := http.NewRequest(http.MethodPost, "/api/bids", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "AuctionID")
	assert.Contains(t, resp.Body.String(), "UserID")
	assert.Contains(t, resp.Body.String(), "Amount")
	assert.Contains(t, resp.Body.String(), "Timestamp")
}
