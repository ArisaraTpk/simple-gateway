package domains

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"simple-gateway/middleware/errors"
)

type VerificationSvc interface {
	Execute(req VerificationReq, ctx *gin.Context, l zerolog.Logger) (*VerificationRes, *errors.APIError)
}

type VerificationReq struct {
	AccountNo string `json:"accountNo" validate:"required,numeric"`
	Amount    int64  `json:"amount" validate:"required,number,gt=0"`
}

type VerificationRes struct {
	TransactionToken string `json:"transactionToken"`
	AccountName      string `json:"accountName"`
}
