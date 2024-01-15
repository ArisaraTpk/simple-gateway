package repository

import (
	"context"
	"github.com/rs/zerolog"
	"simple-gateway/internal/core/ports"
	"simple-gateway/middleware/clientHttp"
)

type notification struct {
	client *clientHttp.Client
}

func NewNotification(client *clientHttp.Client) ports.Notification {
	return &notification{
		client: client,
	}
}

func (r notification) SendNotification(req ports.SendNotificationReq, ctx context.Context, l zerolog.Logger) (*ports.SendNotificationRes, error) {

	cfg := r.client.Config.Apis["noti"]
	var resp ports.SendNotificationRes
	r.client.SetLogger(l)
	err := r.client.Post(cfg.Uri).
		SetContext(ctx).
		SetBody(req).
		Do().
		Into(&resp)

	return &resp, err
}
