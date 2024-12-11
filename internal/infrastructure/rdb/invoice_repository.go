package rdb

import (
	"github.com/take73/invoice-api-example/internal/domain/model"
	"github.com/take73/invoice-api-example/internal/domain/repository"
	"github.com/take73/invoice-api-example/internal/infrastructure/rdb/entity"
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
		PaymentAmount:  entity.BigRat{Rat: invoice.Amount},
		Fee:            entity.BigRat{Rat: invoice.Fee},
		FeeRate:        invoice.FeeRate,
		Tax:            entity.BigRat{Rat: invoice.Tax},
		TaxRate:        invoice.TaxRate,
		TotalAmount:    entity.BigRat{Rat: invoice.TotalAmount},
		DueDate:        invoice.DueDate,
		Status:         string(invoice.Status),
	}

	if err := r.db.Table("invoice").Create(&entity).Error; err != nil {
		return nil, err
	}

	return &model.Invoice{
		ID: entity.ID,
		Organization: &model.Organization{
			ID: entity.OrganizationID,
		},
		Client: &model.Client{
			ID: entity.ClientID,
		},
		IssueDate:   entity.IssueDate,
		Amount:      entity.PaymentAmount.Rat,
		Fee:         entity.Fee.Rat,
		FeeRate:     entity.FeeRate,
		Tax:         entity.Tax.Rat,
		TaxRate:     entity.TaxRate,
		TotalAmount: entity.TotalAmount.Rat,
		DueDate:     entity.DueDate,
		Status:      model.InvoiceStatus(entity.Status),
	}, nil
}
