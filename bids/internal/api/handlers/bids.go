package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/philipjesic/mcg-webapp/bids/internal/api/requests"
	"github.com/philipjesic/mcg-webapp/bids/internal/api/responses"
	"github.com/philipjesic/mcg-webapp/bids/internal/messaging"
	"github.com/philipjesic/mcg-webapp/bids/internal/storage"
)

type Bids struct {
	db        storage.DataStore
	publisher messaging.BidPublisher
}

func CreateBidsHandler(db storage.DataStore, publisher messaging.BidPublisher) *Bids {
	return &Bids{
		db:        db,
		publisher: publisher,
	}
}

func (b *Bids) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	bidID := c.Param("id")
	bid, err := b.db.GetBidByID(ctx, bidID)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Errors: []responses.ErrorMessage{
				{
					Status: http.StatusInternalServerError,
					Title:  "internal server error",
					Detail: "failed to get bid",
				},
			},
		})
		return
	}
	response := createBidResponse([]storage.Bid{bid})
	c.JSON(http.StatusOK, response)
}

func (b *Bids) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	bids, err := b.db.GetBids(ctx, []string{})
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Errors: []responses.ErrorMessage{
				{
					Status: http.StatusInternalServerError,
					Title:  "internal server error",
					Detail: "failed to fetch bids",
				},
			},
		})
		return
	}
	response := createBidResponse(bids)
	c.JSON(http.StatusOK, response)
}

func (b *Bids) Create(c *gin.Context) {
	ctx := c.Request.Context()

	req := requests.BidCreateRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Errors: []responses.ErrorMessage{{
				Status: http.StatusBadRequest,
				Title:  "bad request",
				Detail: err.Error(),
			}},
		})
		return
	}

	bid := storage.Bid{
		ID:        uuid.New().String(),
		AuctionID: req.Data.AuctionID,
		UserID:    req.Data.UserID,
		Amount:    req.Data.Amount,
		Timestamp: req.Data.Timestamp,
	}

	if err := b.db.InsertBid(ctx, bid); err != nil {
		log.Printf("DB error: %v", err)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Errors: []responses.ErrorMessage{{
				Status: http.StatusInternalServerError,
				Title:  "internal server error",
				Detail: "failed to create bid",
			}},
		})
		return
	}

	msg := createBidMessage(bid)

	if err := b.publisher.Publish(messaging.BID_TOPIC, messaging.CREATE_BID, msg); err != nil {
		log.Printf("Publish error: %v", err)
		// Optionally: rollback DB write, or mark for retry, or alert
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Errors: []responses.ErrorMessage{{
				Status: http.StatusInternalServerError,
				Title:  "internal server error",
				Detail: "bid created but failed to publish event",
			}},
		})
		return
	}

	res := createBidResponse([]storage.Bid{bid})
	c.JSON(http.StatusCreated, res)
}

func createBidResponse(bids []storage.Bid) responses.BidsResponse {
	bidResponses := make([]responses.BidResponseBody, 0)
	for _, bid := range bids {
		bidResponses = append(bidResponses, createBidResponseBody(bid))
	}
	return responses.BidsResponse{
		Data: bidResponses,
	}
}

func createBidResponseBody(bid storage.Bid) responses.BidResponseBody {
	return responses.BidResponseBody{
		ID:        bid.ID,
		Type:      "bid",
		AuctionID: bid.AuctionID,
		UserID:    bid.UserID,
		Amount:    bid.Amount,
		Timestamp: bid.Timestamp,
	}
}

func createBidMessage(bid storage.Bid) messaging.BidMessage {
	return messaging.BidMessage{
		Bid: messaging.Bid{
			AuctionID: bid.AuctionID,
			ID:        bid.ID,
			UserID:    bid.UserID,
			Amount:    bid.Amount,
			Timestamp: bid.Timestamp,
		},
	}
}
