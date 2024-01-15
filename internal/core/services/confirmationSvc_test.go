package services_test

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"simple-gateway/internal/adapters/repository"
	"simple-gateway/internal/core/domains"
	"simple-gateway/internal/core/ports"
	"simple-gateway/internal/core/services"
	"simple-gateway/middleware/errors"
	"testing"
	"time"
)

type ConfirmationTestSuite struct {
	suite.Suite
	service   domains.ConfirmationSvc
	transRepo repository.MockTransactionRepo
	prompay   repository.MockPromPay
	noti      repository.MockNotification
}

func (s *ConfirmationTestSuite) SetupSuite() {
	s.prompay = repository.MockPromPay{}
	s.transRepo = repository.MockTransactionRepo{}
	s.noti = repository.MockNotification{}
	validate := validator.New()
	s.service = services.NewConfirmationSvc(&s.prompay, validate, &s.transRepo, &s.noti)
}

func TestConfirmationTestSuite(t *testing.T) {
	suite.Run(t, new(ConfirmationTestSuite))
}

func (suite *ConfirmationTestSuite) TestValidateFailed() {
	req := domains.ConfirmationReq{}
	ctx := gin.Context{}
	l := zerolog.Logger{}
	_, err := suite.service.Execute(req, &ctx, l)
	suite.Equal(err, errors.ErrBadRequest)
}

func (suite *ConfirmationTestSuite) TestSuccess() {
	req := domains.ConfirmationReq{TransactionToken: "jjj"}
	ctx := gin.Context{}
	l := zerolog.Logger{}
	transaction := &ports.TransactionEntity{
		TransactionToken: "jjj",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		AccountName:      "jane",
		AccountNo:        "123456789012",
		PromPayId:        "0666666666",
		Amount:           1,
	}
	prompay := ports.ConfirmationRes{Code: "00"}
	noti := ports.SendNotificationRes{Code: "00"}
	suite.transRepo = repository.MockTransactionRepo{}
	suite.transRepo.On("FindTransaction", mock.AnythingOfType("string")).Return(transaction, nil)
	suite.prompay = repository.MockPromPay{}
	suite.prompay.On("Confirmation", mock.AnythingOfType("ports.ConfirmationReq"), mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("zerolog.Logger")).Return(&prompay, nil)

	suite.noti.On("SendNotification", mock.AnythingOfType("ports.SendNotificationReq"), mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("zerolog.Logger")).Return(&noti, nil)
	res, err := suite.service.Execute(req, &ctx, l)
	suite.Equal(err, (*errors.APIError)(nil))
	suite.Equal(res, &domains.ConfirmationRes{
		Message: "Success",
	})
}

func (suite *ConfirmationTestSuite) TestDbNotFound() {
	req := domains.ConfirmationReq{TransactionToken: "jjj"}
	ctx := gin.Context{}
	l := zerolog.Logger{}
	suite.transRepo = repository.MockTransactionRepo{}
	suite.transRepo.On("FindTransaction", mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound)
	_, err := suite.service.Execute(req, &ctx, l)
	suite.Equal(err, errors.ErrNotFound)
}

func (suite *ConfirmationTestSuite) TestDbFailed() {
	req := domains.ConfirmationReq{TransactionToken: "jjj"}
	ctx := gin.Context{}
	l := zerolog.Logger{}
	suite.transRepo.On("FindTransaction", mock.AnythingOfType("string")).Return(nil, gorm.ErrInvalidDB)
	_, err := suite.service.Execute(req, &ctx, l)
	suite.Equal(err, errors.ErrTechnical)
}

func (suite *ConfirmationTestSuite) TestPrompayFailed() {
	req := domains.ConfirmationReq{TransactionToken: "jjj"}
	ctx := gin.Context{}
	l := zerolog.Logger{}
	transaction := &ports.TransactionEntity{
		TransactionToken: "jjj",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		AccountName:      "jane",
		AccountNo:        "123456789012",
		PromPayId:        "0666666666",
		Amount:           1,
	}
	suite.transRepo = repository.MockTransactionRepo{}
	suite.transRepo.On("FindTransaction", mock.AnythingOfType("string")).Return(transaction, gorm.ErrInvalidDB)
	suite.prompay = repository.MockPromPay{}
	suite.prompay.On("Confirmation", mock.AnythingOfType("ports.ConfirmationReq"), mock.AnythingOfType("context.Context"), mock.AnythingOfType("zerolog.Logger")).Return(nil, new(error))
	_, err := suite.service.Execute(req, &ctx, l)
	suite.Equal(err, errors.ErrTechnical)
}

func (suite *ConfirmationTestSuite) TestPrompayStatusFailed() {
	req := domains.ConfirmationReq{TransactionToken: "jjj"}
	ctx := gin.Context{}
	l := zerolog.Logger{}
	transaction := &ports.TransactionEntity{
		TransactionToken: "jjj",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		AccountName:      "jane",
		AccountNo:        "123456789012",
		PromPayId:        "0666666666",
		Amount:           1,
	}
	prompay := ports.ConfirmationRes{Code: "10"}
	suite.transRepo = repository.MockTransactionRepo{}
	suite.transRepo.On("FindTransaction", mock.AnythingOfType("string")).Return(transaction, nil)
	suite.prompay = repository.MockPromPay{}
	suite.prompay.On("Confirmation", mock.AnythingOfType("ports.ConfirmationReq"), mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("zerolog.Logger")).Return(&prompay, nil)
	_, err := suite.service.Execute(req, &ctx, l)
	suite.Equal(err, errors.ErrBusiness)
}
