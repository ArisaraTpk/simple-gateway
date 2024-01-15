package ports

import (
	"time"
)

type TransactionRepo interface {
	CreateTransaction(TransactionEntity) error
	FindTransaction(transactionToken string) (*TransactionEntity, error)
}

type TransactionEntity struct {
	TransactionToken string    `gorm:"column:transactionToken"`
	CreatedAt        time.Time `gorm:"column:createdAt"`
	UpdatedAt        time.Time `gorm:"column:updatedAt"`
	AccountName      string    `gorm:"column:accountName"`
	AccountNo        string    `gorm:"column:accountNo"`
	PromPayId        string    `gorm:"column:promppayId"`
	Amount           int64     `gorm:"column:amount"`
}
