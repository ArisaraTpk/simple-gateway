package ports

import (
	"context"
	"github.com/rs/zerolog"
)

type PromPay interface {
	Verification(req VerificationReq, ctx context.Context, l zerolog.Logger) (*VerificationRes, error)
	Confirmation(req ConfirmationReq, ctx context.Context, l zerolog.Logger) (*ConfirmationRes, error)
}

type VerificationReq struct {
	AccountNo string `json:"accountNo"`
	PromPayId string `json:"promPayId"`
}

type VerificationRes struct {
	AccountName      string `json:"accountName"`
	PromPayId        string `json:"promPayId"`
	AccountNo        string `json:"accountNo"`
	TransactionToken string `json:"transactionToken"`
}

type ConfirmationReq struct {
	TransactionToken string `json:"transactionToken"`
	PromPayId        string `json:"promPayId"`
	Amount           int64  `json:"amount"`
}

type ConfirmationRes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
