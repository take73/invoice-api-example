package rdb

import (
	"errors"

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
