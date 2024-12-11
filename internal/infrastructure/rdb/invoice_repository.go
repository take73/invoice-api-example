package rdb

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"github.com/take73/invoice-api-example/internal/domain/model"
	"github.com/take73/invoice-api-example/internal/domain/repository"
	"github.com/take73/invoice-api-example/internal/infrastructure/rdb/entity"
	"github.com/take73/invoice-api-example/internal/shared/validation"
	"gorm.io/gorm"
)

type InvoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) repository.Invoice {
	return &InvoiceRepository{db: db}
}

// Create は請求書をデータベースに保存
func (r *InvoiceRepository) Create(invoice *model.Invoice) (*model.Invoice, error) {
	entity := entity.Invoice{
		OrganizationID: invoice.Organization.ID,
		ClientID:       invoice.Client.ID,
		IssueDate:      invoice.IssueDate,
		PaymentAmount:  invoice.Amount,
		Fee:            invoice.Fee,
		FeeRate:        decimal.NewFromFloat(invoice.FeeRate),
		Tax:            invoice.Tax,
		TaxRate:        decimal.NewFromFloat(invoice.TaxRate),
		TotalAmount:    invoice.TotalAmount,
		DueDate:        invoice.DueDate,
		Status:         string(invoice.Status),
	}

	if err := r.db.Table("invoice").Create(&entity).Error; err != nil {
		return nil, err
	}

	taxRate, _ := entity.TaxRate.Float64()
	if validation.ValidRate(taxRate) {
		return nil, errors.New("invalid taxRate: must be between 0.0 and 1.0")
	}
	feeRate, _ := entity.FeeRate.Float64()
	if validation.ValidRate(feeRate) {
		return nil, errors.New("invalid feeRate: must be between 0.0 and 1.0")
	}

	// domain modelに変換して返す
	return &model.Invoice{
		ID: entity.ID,
		Organization: &model.Organization{
			ID:   entity.OrganizationID,
			Name: invoice.Organization.Name,
		},
		Client: &model.Client{
			ID:   entity.ClientID,
			Name: invoice.Client.Name,
		},
		IssueDate:   entity.IssueDate,
		Amount:      entity.PaymentAmount,
		Fee:         entity.Fee,
		FeeRate:     feeRate,
		Tax:         entity.Tax,
		TaxRate:     taxRate,
		TotalAmount: entity.TotalAmount,
		DueDate:     entity.DueDate,
		Status:      model.InvoiceStatus(entity.Status),
	}, nil
}

func (r *InvoiceRepository) FindByDueDateRange(startDate, endDate time.Time) ([]*model.Invoice, error) {
	var entities []entity.Invoice

	err := r.db.Table("invoice").
		Where("due_date >= ? AND due_date <= ?", startDate, endDate).
		Order("due_date asc").
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	// ドメインモデルに変換
	invoices := make([]*model.Invoice, len(entities))
	for i, e := range entities {
		taxRate, _ := e.TaxRate.Float64()
		feeRate, _ := e.FeeRate.Float64()

		invoices[i] = &model.Invoice{
			ID: e.ID,
			Organization: &model.Organization{
				ID: e.OrganizationID,
			},
			Client: &model.Client{
				ID: e.ClientID,
			},
			IssueDate:   e.IssueDate,
			Amount:      e.PaymentAmount,
			Fee:         e.Fee,
			FeeRate:     feeRate,
			Tax:         e.Tax,
			TaxRate:     taxRate,
			TotalAmount: e.TotalAmount,
			DueDate:     e.DueDate,
			Status:      model.InvoiceStatus(e.Status),
		}
	}

	return invoices, nil
}
