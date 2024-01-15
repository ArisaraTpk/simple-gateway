package repository

import (
	"context"
	"github.com/rs/zerolog"
	"simple-gateway/internal/core/ports"
	"simple-gateway/middleware/clientHttp"
)

type promPay struct {
	client *clientHttp.Client
}

func NewPromPay(client *clientHttp.Client) ports.PromPay {
	return &promPay{
		client: client,
	}
}

func (r promPay) Verification(req ports.VerificationReq, ctx context.Context, l zerolog.Logger) (*ports.VerificationRes, error) {

	cfg := r.client.Config.Apis["verification"]
	var resp ports.VerificationRes
	r.client.SetLogger(l)
	err := r.client.Post(cfg.Uri).
		SetContext(ctx).
		SetBody(req).
		Do().
		Into(&resp)

	return &resp, err
}

func (r promPay) Confirmation(req ports.ConfirmationReq, ctx context.Context, l zerolog.Logger) (*ports.ConfirmationRes, error) {
	cfg := r.client.Config.Apis["confirmation"]
	var resp ports.ConfirmationRes
	r.client.SetLogger(l)
	err := r.client.Post(cfg.Uri).
		SetContext(ctx).
		SetBody(req).
		Do().
		Into(&resp)

	return &resp, err
}
