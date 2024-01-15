package domains

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"simple-gateway/middleware/errors"
)

type ConfirmationSvc interface {
	Execute(req ConfirmationReq, ctx *gin.Context, l zerolog.Logger) (*ConfirmationRes, *errors.APIError)
}

type ConfirmationReq struct {
	TransactionToken string `json:"transactionToken" validate:"required"`
}

type ConfirmationRes struct {
	Message string `json:"message"`
}
