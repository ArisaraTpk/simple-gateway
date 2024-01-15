package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"simple-gateway/internal/core/domains"
)

type VerificationHdl interface {
	Verification(c *gin.Context)
}

type verificationHdl struct {
	svc domains.VerificationSvc
}

func NewVerificationHdl(svc domains.VerificationSvc) VerificationHdl {
	return &verificationHdl{
		svc: svc,
	}
}

func (h verificationHdl) Verification(c *gin.Context) {
	l := zerolog.New(os.Stdout).With().Timestamp().Logger()

	var req domains.VerificationReq
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
