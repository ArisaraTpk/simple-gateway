package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog"
	"simple-gateway/internal/core/domains"
	"simple-gateway/internal/core/ports"
	"simple-gateway/middleware/errors"
	"simple-gateway/utils"
	"time"
)

type verificationSvc struct {
	prompay         ports.PromPay
	validator       *validator.Validate
	transactionRepo ports.TransactionRepo
}

func NewVerificationSvc(prompay ports.PromPay, validator *validator.Validate, transactionRepo ports.TransactionRepo) domains.VerificationSvc {
	return &verificationSvc{
		prompay:         prompay,
		validator:       validator,
		transactionRepo: transactionRepo,
	}
}

func (s verificationSvc) Execute(req domains.VerificationReq, ctx *gin.Context, l zerolog.Logger) (*domains.VerificationRes, *errors.APIError) {

	if err := s.validate(req, l); err != nil {
		return nil, err
	}

	reqVer := s.buildVerificationReq(req)
	res, err := s.prompay.Verification(reqVer, ctx.Copy(), l)
	if err != nil {
		l.Error().Err(err).Msg(fmt.Sprintf("Error call prompay.Verification %v", err.Error()))
		return nil, errors.ErrTechnical
	}

	record := s.buildTransaction(res, req)
	err = s.transactionRepo.CreateTransaction(record)
	if err != nil {
		l.Error().Err(err).Msg(fmt.Sprintf("Error create record  %v", err.Error()))
		return nil, errors.ErrTechnical
	}

	return s.buildVerificationRes(res), nil
}

func (s verificationSvc) validate(req domains.VerificationReq, l zerolog.Logger) *errors.APIError {
	if err := s.validator.Struct(req); err != nil {
		l.Error().Err(err).Msg(fmt.Sprintf("validate errors %v", err))
		return errors.ErrBadRequest
	}
	return nil
}

func (s verificationSvc) buildVerificationReq(req domains.VerificationReq) ports.VerificationReq {
	data := ports.VerificationReq{}

	if utils.ValidateMobile(req.AccountNo) {
		data.PromPayId = req.AccountNo
	} else {
		data.AccountNo = req.AccountNo
	}
	return data
}

func (s verificationSvc) buildVerificationRes(res *ports.VerificationRes) *domains.VerificationRes {
	return &domains.VerificationRes{
		TransactionToken: res.TransactionToken,
		AccountName:      res.AccountName}
}

func (s verificationSvc) buildTransaction(data *ports.VerificationRes, req domains.VerificationReq) ports.TransactionEntity {
	t := time.Now()
	return ports.TransactionEntity{
		TransactionToken: data.TransactionToken,
		AccountNo:        data.AccountNo,
		PromPayId:        data.PromPayId,
		AccountName:      data.AccountName,
		CreatedAt:        t,
		UpdatedAt:        t,
		Amount:           req.Amount,
	}
}
