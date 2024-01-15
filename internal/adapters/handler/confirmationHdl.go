package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"simple-gateway/internal/core/domains"
)

type ConfirmationHdl interface {
	Confirmation(c *gin.Context)
}

type confirmationHdl struct {
	svc domains.ConfirmationSvc
}

func NewConfirmationHdl(svc domains.ConfirmationSvc) ConfirmationHdl {
	return &confirmationHdl{
		svc: svc,
	}
}

func (h confirmationHdl) Confirmation(c *gin.Context) {
	l := zerolog.New(os.Stdout).With().Timestamp().Logger()

	var req domains.ConfirmationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err.Error(),
		})
		return
	}
	res, err := h.svc.Execute(req, c, l)
	if err != nil {
		c.JSON(err.GetCode(), gin.H{
			"errors": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
