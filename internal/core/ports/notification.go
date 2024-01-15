package ports

import (
	"context"
	"github.com/rs/zerolog"
)

type Notification interface {
	SendNotification(req SendNotificationReq, ctx context.Context, l zerolog.Logger) (*SendNotificationRes, error)
}

type SendNotificationReq struct {
	AccountName      string `json:"accountName"`
	Amount           int64  `json:"amount"`
	PromPayId        string `json:"promPayId"`
	AccountNo        string `json:"accountNo"`
	TransactionToken string `json:"transactionToken"`
}

type SendNotificationRes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
