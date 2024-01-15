package repository

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"simple-gateway/internal/core/ports"
)

type MockPromPay struct {
	mock.Mock
}

func NewMockPromPay() ports.PromPay {
	return &MockPromPay{}
}

func (r MockPromPay) Verification(req ports.VerificationReq, ctx context.Context, l zerolog.Logger) (*ports.VerificationRes, error) {
	args := r.Called(req, ctx, l)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.VerificationRes), args.Error(1)
}

func (r MockPromPay) Confirmation(req ports.ConfirmationReq, ctx context.Context, l zerolog.Logger) (*ports.ConfirmationRes, error) {
	args := r.Called(req, ctx, l)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.ConfirmationRes), args.Error(1)
}
