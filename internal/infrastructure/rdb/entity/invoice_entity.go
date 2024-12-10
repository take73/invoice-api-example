package entity

import "time"

// InvoiceEntity はGORMやデータベース関連の詳細を含むモデル
type InvoiceEntity struct {
	ID             int       `gorm:"primaryKey;autoIncrement"`
	OrganizationID int       `gorm:"column:organization_id"`
	ClientID       int       `gorm:"column:client_id"`
	IssueDate      time.Time `gorm:"column:issue_date"`
	PaymentAmount  float64   `gorm:"column:payment_amount"`
	Fee            float64   `gorm:"column:fee"`
	FeeRate        float64   `gorm:"column:fee_rate"`
	Tax            float64   `gorm:"column:tax"`
	TaxRate        float64   `gorm:"column:tax_rate"`
	TotalAmount    float64   `gorm:"column:total_amount"`
	DueDate        time.Time `gorm:"column:due_date"`
	Status         string    `gorm:"type:enum('pending','processing','paid','error');default:'pending'"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
