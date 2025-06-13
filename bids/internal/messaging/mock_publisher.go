package messaging

import (
	"github.com/stretchr/testify/mock"
)

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(topic, key string, bidMsg BidMessage) error {
	args := m.Called(topic, key, bidMsg)
	return args.Error(0)
}
