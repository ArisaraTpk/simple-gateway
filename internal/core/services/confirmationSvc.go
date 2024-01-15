package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"simple-gateway/internal/core/domains"
	"simple-gateway/internal/core/ports"
	"simple-gateway/middleware/errors"
)

type confirmationSvc struct {
	prompay         ports.PromPay
	validator       *validator.Validate
	transactionRepo ports.TransactionRepo
	noti            ports.Notification
}

func NewConfirmationSvc(prompay ports.PromPay, validator *validator.Validate, transactionRepo ports.TransactionRepo, noti ports.Notification) domains.ConfirmationSvc {
	return &confirmationSvc{
		prompay:         prompay,
		validator:       validator,
		transactionRepo: transactionRepo,
		noti:            noti,
	}
}

func (s confirmationSvc) Execute(req domains.ConfirmationReq, ctx *gin.Context, l zerolog.Logger) (*domains.ConfirmationRes, *errors.APIError) {
	if err := s.validate(req, l); err != nil {
		return nil, err
	}

	data, err := s.transactionRepo.FindTransaction(req.TransactionToken)
	if err != nil {
		l.Error().Err(err).Msg(fmt.Sprintf("FindTransaction errors %v", err))
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrTechnical
	}

	errConfirm := s.Confirm(data, ctx, l)
	if errConfirm != nil {
		return nil, errConfirm
	}

	reqNoti := s.buildSendNotiReq(data)
	_, err = s.noti.SendNotification(reqNoti, ctx.Copy(), l)
	if err != nil {
		l.Error().Err(err).Msg(fmt.Sprintf("noti.SendNotification errors %v", err))
		return nil, errors.ErrTechnical
	}

	return s.buildConfirmationRes(), nil
}

func (s confirmationSvc) validate(req domains.ConfirmationReq, l zerolog.Logger) *errors.APIError {
	if err := s.validator.Struct(req); err != nil {
		l.Error().
			Err(err).
			Msg(fmt.Sprintf("validate errors %v", err))
		return errors.ErrBadRequest
	}
	return nil
}

func (s confirmationSvc) buildConfirmationReq(req *ports.TransactionEntity) ports.ConfirmationReq {
	return ports.ConfirmationReq{
		TransactionToken: req.TransactionToken,
		Amount:           req.Amount,
		PromPayId:        req.PromPayId,
	}
}

func (s confirmationSvc) buildSendNotiReq(data *ports.TransactionEntity) ports.SendNotificationReq {
	return ports.SendNotificationReq{
		TransactionToken: data.TransactionToken,
		AccountName:      data.AccountName,
		AccountNo:        data.AccountNo,
		PromPayId:        data.PromPayId,
		Amount:           data.Amount,
	}
}

func (s confirmationSvc) buildConfirmationRes() *domains.ConfirmationRes {
	return &domains.ConfirmationRes{
		Message: "Success",
	}
}

func (s confirmationSvc) Confirm(data *ports.TransactionEntity, ctx *gin.Context, l zerolog.Logger) *errors.APIError {
	reqCon := s.buildConfirmationReq(data)
	res, err := s.prompay.Confirmation(reqCon, ctx.Copy(), l)
	if err != nil {
		l.Error().
			Err(err).
			Msg(fmt.Sprintf("prompay.Confirmation errors %v", err))
		return errors.ErrTechnical
	}
	if res.Code != "00" {
		l.Error().Err(err).Msg(fmt.Sprintf("prompay.Confirmation ResponseCode %v", res.Code))
		return errors.ErrBusiness
	}

	return nil
}
