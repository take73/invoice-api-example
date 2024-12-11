package rdb

import (
	"github.com/take73/invoice-api-example/internal/domain/model"
	"github.com/take73/invoice-api-example/internal/infrastructure/rdb/entity"
	"gorm.io/gorm"
)

type InvoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

// Create は請求書をデータベースに保存
func (r *InvoiceRepository) Create(invoice model.Invoice) (*model.Invoice, error) {
	// ドメインモデルを永続化モデルに変換
	entity := entity.InvoiceEntity{
		OrganizationID: invoice.Organization.ID,
		ClientID:       invoice.Client.ID,
		IssueDate:      invoice.IssueDate,
		PaymentAmount:  invoice.Amount,
		Fee:            invoice.Fee,
		FeeRate:        invoice.FeeRate,
		Tax:            invoice.Tax,
		TaxRate:        invoice.TaxRate,
		TotalAmount:    invoice.TotalAmount,
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
		Amount:      entity.PaymentAmount,
		Fee:         entity.Fee,
		FeeRate:     entity.FeeRate,
		Tax:         entity.Tax,
		TaxRate:     entity.TaxRate,
		TotalAmount: entity.TotalAmount,
		DueDate:     entity.DueDate,
		Status:      model.InvoiceStatus(entity.Status),
	}, nil
}
