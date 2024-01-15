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
)

type VerificationTestSuite struct {
	suite.Suite
	service   domains.VerificationSvc
	transRepo repository.MockTransactionRepo
	prompay   repository.MockPromPay
}

func (s *VerificationTestSuite) SetupSuite() {
	s.prompay = repository.MockPromPay{}
	s.transRepo = repository.MockTransactionRepo{}
	validate := validator.New()
	s.service = services.NewVerificationSvc(&s.prompay, validate, &s.transRepo)
}

func TestVerificationTestSuite(t *testing.T) {
	suite.Run(t, new(VerificationTestSuite))
}

func (suite *VerificationTestSuite) TestValidateFailedAccountNo() {
	req := domains.VerificationReq{
		AccountNo: "",
		Amount:    1,
	}
	ctx := gin.Context{}
	l := zerolog.Logger{}
	_, err := suite.service.Execute(req, &ctx, l)
	suite.Equal(err, errors.ErrBadRequest)
}
func (suite *VerificationTestSuite) TestValidateFailedAmount() {
	req := domains.VerificationReq{
		AccountNo: "12345",
		Amount:    0,
	}
	ctx := gin.Context{}
	l := zerolog.Logger{}
	_, err := suite.service.Execute(req, &ctx, l)
	suite.Equal(err, errors.ErrBadRequest)
}

func (suite *VerificationTestSuite) TestSuccess() {
	req := domains.VerificationReq{
		AccountNo: "123456789012",
		Amount:    1,
	}
	ctx := gin.Context{}
	l := zerolog.Logger{}
	prompay := ports.VerificationRes{
		AccountName:      "jane",
		PromPayId:        "0666666666",
		AccountNo:        "123456789012",
		TransactionToken: "jjj",
	}
	suite.transRepo = repository.MockTransactionRepo{}
	suite.transRepo.On("CreateTransaction", mock.AnythingOfType("ports.TransactionEntity")).Return(nil)
	suite.prompay = repository.MockPromPay{}
	suite.prompay.On("Verification", mock.AnythingOfType("ports.VerificationReq"), mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("zerolog.Logger")).Return(&prompay, nil)
	res, err := suite.service.Execute(req, &ctx, l)
	suite.Equal(err, (*errors.APIError)(nil))
	suite.Equal(res, &domains.VerificationRes{
		TransactionToken: "jjj",
		AccountName:      "jane",
	})
}
func (suite *VerificationTestSuite) TestPromPayFailed() {
	req := domains.VerificationReq{
		AccountNo: "123456789012",
		Amount:    1}
	ctx := gin.Context{}
	l := zerolog.Logger{}
	suite.transRepo = repository.MockTransactionRepo{}
	suite.transRepo.On("CreateTransaction", mock.AnythingOfType("ports.TransactionEntity")).Return(nil)
	suite.prompay = repository.MockPromPay{}
	suite.prompay.On("Verification", mock.AnythingOfType("ports.VerificationReq"), mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("zerolog.Logger")).Return(nil, errors.ErrTechnical)
	_, err := suite.service.Execute(req, &ctx, l)
	suite.Equal(err, errors.ErrTechnical)
}

func (suite *VerificationTestSuite) TestDbFailed() {
	req := domains.VerificationReq{
		AccountNo: "123456789012",
		Amount:    1}
	ctx := gin.Context{}
	l := zerolog.Logger{}
	prompay := ports.VerificationRes{
		AccountName:      "jane",
		PromPayId:        "0666666666",
		AccountNo:        "123456789012",
		TransactionToken: "jjj",
	}
	suite.transRepo = repository.MockTransactionRepo{}
	suite.transRepo.On("CreateTransaction", mock.AnythingOfType("ports.TransactionEntity")).Return(gorm.ErrInvalidDB)
	suite.prompay = repository.MockPromPay{}
	suite.prompay.On("Verification", mock.AnythingOfType("ports.VerificationReq"), mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("zerolog.Logger")).Return(&prompay, nil)
	_, err := suite.service.Execute(req, &ctx, l)
	suite.Equal(err, errors.ErrTechnical)
}

//
//func (suite *VerificationTestSuite) TestPrompayFailed() {
//	req := domains.VerificationReq{TransactionToken: "jjj"}
//	ctx := gin.Context{}
//	l := zerolog.Logger{}
//	transaction := &ports.TransactionEntity{
//		TransactionToken: "jjj",
//		CreatedAt:        time.Now(),
//		UpdatedAt:        time.Now(),
//		AccountName:      "jane",
//		AccountNo:        "123456789012",
//		PromPayId:        "0666666666",
//		Amount:           1,
//	}
//	suite.transRepo = repository.MockTransactionRepo{}
//	suite.transRepo.On("FindTransaction", mock.AnythingOfType("string")).Return(transaction, gorm.ErrInvalidDB)
//	suite.prompay = repository.MockPromPay{}
//	suite.prompay.On("Verification", mock.AnythingOfType("ports.VerificationReq"), mock.AnythingOfType("context.Context"), mock.AnythingOfType("zerolog.Logger")).Return(nil, new(error))
//	_, err := suite.service.Execute(req, &ctx, l)
//	suite.Equal(err, errors.ErrTechnical)
//}
//
//func (suite *VerificationTestSuite) TestPrompayStatusFailed() {
//	req := domains.VerificationReq{TransactionToken: "jjj"}
//	ctx := gin.Context{}
//	l := zerolog.Logger{}
//	transaction := &ports.TransactionEntity{
//		TransactionToken: "jjj",
//		CreatedAt:        time.Now(),
//		UpdatedAt:        time.Now(),
//		AccountName:      "jane",
//		AccountNo:        "123456789012",
//		PromPayId:        "0666666666",
//		Amount:           1,
//	}
//	prompay := ports.VerificationRes{Code: "10"}
//	suite.transRepo = repository.MockTransactionRepo{}
//	suite.transRepo.On("FindTransaction", mock.AnythingOfType("string")).Return(transaction, nil)
//	suite.prompay = repository.MockPromPay{}
//	suite.prompay.On("Verification", mock.AnythingOfType("ports.VerificationReq"), mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("zerolog.Logger")).Return(&prompay, nil)
//	_, err := suite.service.Execute(req, &ctx, l)
//	suite.Equal(err, errors.ErrBusiness)
//}
