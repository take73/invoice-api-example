package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

// Invoice ORMのEntity
type Invoice struct {
	ID             uint            `gorm:"primaryKey;autoIncrement;column:invoice_id"`
	OrganizationID uint            `gorm:"column:organization_id;not null"`
	ClientID       uint            `gorm:"column:client_id;not null"`
	IssueDate      time.Time       `gorm:"column:issue_date;not null"`
	PaymentAmount  decimal.Decimal `gorm:"column:payment_amount;type:decimal(10,2);not null"`
	Fee            decimal.Decimal `gorm:"column:fee;type:decimal(10,2)"`
	FeeRate        decimal.Decimal `gorm:"column:fee_rate;type:decimal(5,2)"`
	Tax            decimal.Decimal `gorm:"column:tax;type:decimal(10,2)"`
	TaxRate        decimal.Decimal `gorm:"column:tax_rate;type:decimal(5,2)"`
	TotalAmount    decimal.Decimal `gorm:"column:total_amount;type:decimal(10,2);not null"`
	DueDate        time.Time       `gorm:"column:due_date;not null"`
	Status         string          `gorm:"column:status;type:enum('pending','processing','paid','error');default:'pending'"`
	CreatedAt      time.Time       `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time       `gorm:"column:updated_at;autoUpdateTime"`

	// Associations
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Client       Client       `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name used by GORM.
func (Invoice) TableName() string {
	return "invoice"
}
