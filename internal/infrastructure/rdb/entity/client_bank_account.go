package entity

import "time"

// ClientBankAccount ORMのEntity
type ClientBankAccount struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;column:account_id"`
	ClientID      uint      `gorm:"column:client_id;not null"`
	BankName      string    `gorm:"column:bank_name;not null"`
	BranchName    string    `gorm:"column:branch_name;not null"`
	AccountNumber string    `gorm:"column:account_number;not null"`
	AccountName   string    `gorm:"column:account_name;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Associations
	Client Client `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name used by GORM.
func (ClientBankAccount) TableName() string {
	return "client_bank_account"
}
