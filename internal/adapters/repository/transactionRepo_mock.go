package repository

import (
	"github.com/stretchr/testify/mock"
	"simple-gateway/internal/core/ports"
)

type MockTransactionRepo struct {
	mock.Mock
}

func NewMockTransactionRepo() ports.TransactionRepo {
	return &MockTransactionRepo{}
}

func (r MockTransactionRepo) CreateTransaction(data ports.TransactionEntity) error {
	args := r.Called(data)
	return args.Error(0)
}

func (r MockTransactionRepo) FindTransaction(transactionToken string) (*ports.TransactionEntity, error) {
	args := r.Called(transactionToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.TransactionEntity), args.Error(1)
}
