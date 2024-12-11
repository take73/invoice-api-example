package entity

import "time"

// TaxRate ORMのEntity
type TaxRate struct {
	ID        uint       `gorm:"primaryKey;autoIncrement;column:tax_rate_id"`
	StartDate time.Time  `gorm:"column:start_date;not null"` // 税率の適用開始日
	EndDate   *time.Time `gorm:"column:end_date"`            // 税率の適用終了日（NULLなら現在も有効）
	Rate      float64    `gorm:"column:rate;not null"`       // 税率（例: 10.00 = 10%）
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (TaxRate) TableName() string {
	return "tax_rate"
}
