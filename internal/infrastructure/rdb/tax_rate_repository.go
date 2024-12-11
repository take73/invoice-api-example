package rdb

import (
	"errors"
	"fmt"
	"time"

	"github.com/take73/invoice-api-example/internal/domain/repository"
	"github.com/take73/invoice-api-example/internal/infrastructure/rdb/entity"
	"gorm.io/gorm"
)

type TaxRateRepository struct {
	db *gorm.DB
}

func NewTaxRateRepository(db *gorm.DB) repository.TaxRate {
	return &TaxRateRepository{db: db}
}

// GetRateByDate 指定した日付に適用される税率を取得します
func (r *TaxRateRepository) GetRateByDate(date time.Time) (float64, error) {
	var taxRate entity.TaxRate

	if err := r.db.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", date, date).
		Order("start_date DESC").
		First(&taxRate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("no tax rate found for date %s: %w", date.Format("2006-01-02"), err)
		}
		return 0, fmt.Errorf("failed to retrieve tax rate for date %s: %w", date.Format("2006-01-02"), err)
	}

	return taxRate.Rate, nil
}
