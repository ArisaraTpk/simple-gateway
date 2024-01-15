package repository

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"simple-gateway/internal/core/ports"
)

type MockNotification struct {
	mock.Mock
}

func NewMockNotification() ports.Notification {
	return &notification{}
}

func (r MockNotification) SendNotification(req ports.SendNotificationReq, ctx context.Context, l zerolog.Logger) (*ports.SendNotificationRes, error) {
	args := r.Called(req, ctx, l)

	return args.Get(0).(*ports.SendNotificationRes), args.Error(1)
}
