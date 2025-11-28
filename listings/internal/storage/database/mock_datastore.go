package database

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockDataStore struct {
	mock.Mock
}

func (m *MockDataStore) GetListingByID(ctx context.Context, id string) (Listing, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Listing), args.Error(1)
}

func (m *MockDataStore) GetListings(ctx context.Context, ids []string) ([]Listing, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]Listing), args.Error(1)
}

func (m *MockDataStore) InsertListing(ctx context.Context, listing Listing) error {
	args := m.Called(ctx, listing)
	return args.Error(0)
}

func (m *MockDataStore) UpdateListingBid(ctx context.Context, auctionID, userID string, amount int, timestamp time.Time) error {
	args := m.Called(ctx, auctionID, userID, amount, timestamp)
	return args.Error(0)
}