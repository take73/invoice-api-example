package entity

import "time"

// Invoice ORMのEntity
type Invoice struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;column:invoice_id"`
	OrganizationID uint      `gorm:"column:organization_id;not null"`
	ClientID       uint      `gorm:"column:client_id;not null"`
	IssueDate      time.Time `gorm:"column:issue_date;not null"`
	PaymentAmount  BigRat    `gorm:"column:payment_amount;not null"`
	Fee            BigRat    `gorm:"column:fee"`
	FeeRate        float64   `gorm:"column:fee_rate"`
	Tax            BigRat    `gorm:"column:tax"`
	TaxRate        float64   `gorm:"column:tax_rate"`
	TotalAmount    BigRat    `gorm:"column:total_amount;not null"`
	DueDate        time.Time `gorm:"column:due_date;not null"`
	Status         string    `gorm:"column:status;type:enum('pending','processing','paid','error');default:'pending'"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Associations
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Client       Client       `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name used by GORM.
func (Invoice) TableName() string {
	return "invoice"
}
