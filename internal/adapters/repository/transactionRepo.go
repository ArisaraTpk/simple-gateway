package repository

import (
	"gorm.io/gorm"
	"simple-gateway/internal/core/ports"
)

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) ports.TransactionRepo {
	return &transactionRepo{
		db: db,
	}
}

func (r transactionRepo) CreateTransaction(data ports.TransactionEntity) error {
	res := r.db.Create(&data)
	return res.Error
}

func (r transactionRepo) FindTransaction(transactionToken string) (*ports.TransactionEntity, error) {
	var result ports.TransactionEntity
	res := r.db.Where("transactionToken = ?", transactionToken).First(&result)
	return &result, res.Error
}
