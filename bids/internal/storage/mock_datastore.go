package storage

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockDataStore struct {
	mock.Mock
}

func (m *MockDataStore) GetBidByID(ctx context.Context, id string) (Bid, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Bid), args.Error(1)
}

func (m *MockDataStore) GetBids(ctx context.Context, ids []string) ([]Bid, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]Bid), args.Error(1)
}

func (m *MockDataStore) InsertBid(ctx context.Context, b Bid) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}