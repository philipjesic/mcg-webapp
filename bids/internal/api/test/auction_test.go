package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/bids/internal/api/handlers"
	"github.com/philipjesic/mcg-webapp/bids/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

type mockAuctionService struct {
	subscribeFn        func(string) chan service.BidEvent
	unsubscribeFn      func(string, chan service.BidEvent)
	broadcastBidEventFn func(context.Context, amqp.Delivery) error
}

func (m *mockAuctionService) Subscribe(auctionID string) chan service.BidEvent {
	if m.subscribeFn != nil {
		return m.subscribeFn(auctionID)
	}
	return nil
}

func (m *mockAuctionService) Unsubscribe(auctionID string, bidChannel chan service.BidEvent) {
	if m.unsubscribeFn != nil {
		m.unsubscribeFn(auctionID, bidChannel)
	}
}

func (m *mockAuctionService) BroadcastBidEvent(ctx context.Context, msg amqp.Delivery) error {
	if m.broadcastBidEventFn != nil {
		return m.broadcastBidEventFn(ctx, msg)
	}
	return nil
}

type MockResponseWriter struct {
    *httptest.ResponseRecorder
}

func (m *MockResponseWriter) CloseNotify() <-chan bool {
    return make(chan bool) // Return a dummy channel
}

func TestAuctionStream_SubscribesWithRouteParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var subscribedAuctionID string
	ch := make(chan service.BidEvent)
	close(ch)

	svc := &mockAuctionService{
		subscribeFn: func(auctionID string) chan service.BidEvent {
			subscribedAuctionID = auctionID
			return ch
		},
		unsubscribeFn: func(string, chan service.BidEvent) {},
	}

	handler := handlers.CreateAuctionHandler(svc)

	r := gin.New()
	r.GET("/api/streams/auctions/:id", handler.Stream)

	req := httptest.NewRequest(http.MethodGet, "/api/streams/auctions/auction-123", nil)
	w := httptest.NewRecorder()
	mockW := &MockResponseWriter{w}

	r.ServeHTTP(mockW, req)

	require.Equal(t, "auction-123", subscribedAuctionID)
}

func TestAuctionStream_SendsConnectedAndBidEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	bidCh := make(chan service.BidEvent, 1)

	svc := &mockAuctionService{
		subscribeFn: func(string) chan service.BidEvent {
			return bidCh
		},
		unsubscribeFn: func(string, chan service.BidEvent) {},
	}

	handler := handlers.CreateAuctionHandler(svc)

	r := gin.New()
	r.GET("/api/streams/auctions/:id", handler.Stream)

	req := httptest.NewRequest(http.MethodGet, "/api/streams/auctions/auction-123", nil)
	w := httptest.NewRecorder()
	mockW := &MockResponseWriter{w}

	done := make(chan struct{})
	go func() {
		r.ServeHTTP(mockW, req)
		close(done)
	}()

	bidCh <- service.BidEvent{
		Data: service.BidEventData{
			AuctionID: "auction-123",
			ID: "test-id",
			UserID: "user-1",
			Amount: 42000,
			Timestamp: time.Now(),
		},
		Event: "bid.create",
	}
	close(bidCh)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("stream handler did not finish")
	}

	body := w.Body.String()

	// Need to serialize json to object then check field.

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, body, "event:Connected")
	require.Contains(t, body, "data:Listening for auction bids...")
	require.Contains(t, body, "event:Bid")
	require.Contains(t, body, "auction-123")
	require.Contains(t, body, "42000")
	require.Contains(t, body, "user-1")
}

func TestAuctionStream_SendsDisconnectWhenChannelClosed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	bidCh := make(chan service.BidEvent)

	svc := &mockAuctionService{
		subscribeFn: func(string) chan service.BidEvent {
			return bidCh
		},
		unsubscribeFn: func(string, chan service.BidEvent) {},
	}

	handler := handlers.CreateAuctionHandler(svc)

	r := gin.New()
	r.GET("/api/streams/auctions/:id", handler.Stream)

	req := httptest.NewRequest(http.MethodGet, "/api/streams/auctions/auction-999", nil)
	w := httptest.NewRecorder()
	mockW := &MockResponseWriter{w}

	done := make(chan struct{})
	go func() {
		r.ServeHTTP(mockW, req)
		close(done)
	}()

	close(bidCh)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("stream handler did not finish")
	}

	body := w.Body.String()

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, body, "event:Connected")
	require.Contains(t, body, "event:Disconnect")
	require.Contains(t, body, "data:The auction is now closed")
}

func TestAuctionStream_CallsUnsubscribeOnExit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	bidCh := make(chan service.BidEvent)

	var mu sync.Mutex
	var unsubCalled bool
	var unsubAuctionID string
	var unsubChannel chan service.BidEvent

	svc := &mockAuctionService{
		subscribeFn: func(string) chan service.BidEvent {
			return bidCh
		},
		unsubscribeFn: func(auctionID string, ch chan service.BidEvent) {
			mu.Lock()
			defer mu.Unlock()
			unsubCalled = true
			unsubAuctionID = auctionID
			unsubChannel = ch
		},
	}

	handler := handlers.CreateAuctionHandler(svc)

	r := gin.New()
	r.GET("/api/streams/auctions/:id", handler.Stream)

	req := httptest.NewRequest(http.MethodGet, "/api/streams/auctions/auction-777", nil)
	w := httptest.NewRecorder()
	mockW := &MockResponseWriter{w}

	done := make(chan struct{})
	go func() {
		r.ServeHTTP(mockW, req)
		close(done)
	}()

	close(bidCh)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("stream handler did not finish")
	}

	mu.Lock()
	defer mu.Unlock()

	require.True(t, unsubCalled)
	require.Equal(t, "auction-777", unsubAuctionID)
	require.Equal(t, bidCh, unsubChannel)
}