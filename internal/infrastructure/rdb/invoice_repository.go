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
	var createdInvoice *model.Invoice

	// トランザクションの開始（複数のリポジトリをまたぐ管理をしたい場合は、ドメインサービスを作り、そこでトランザクションを管理する）
	err := r.db.Transaction(func(tx *gorm.DB) error {
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

		// データベースに登録
		if err := tx.Table("invoice").Create(&entity).Error; err != nil {
			return err
		}

		// 税率・手数料率の妥当性を検証
		taxRate, _ := entity.TaxRate.Float64()
		if !validation.ValidRate(taxRate) {
			return errors.New("invalid taxRate: must be between 0.0 and 1.0")
		}
		feeRate, _ := entity.FeeRate.Float64()
		if !validation.ValidRate(feeRate) {
			return errors.New("invalid feeRate: must be between 0.0 and 1.0")
		}

		// ドメインモデルに変換
		createdInvoice = &model.Invoice{
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
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdInvoice, nil
}

func (r *InvoiceRepository) FindByDueDateRange(startDate, endDate time.Time) ([]*model.Invoice, error) {
	type listInvoiceItem struct {
		entity.Invoice
		OrganizationName string `gorm:"column:organization_name"`
		ClientName       string `gorm:"column:client_name"`
	}

	var entities []listInvoiceItem

	err := r.db.Table("invoice").
		Select("invoice.*, organization.name AS organization_name, client.name AS client_name").
		Joins("JOIN organization ON invoice.organization_id = organization.organization_id").
		Joins("JOIN client ON invoice.client_id = client.client_id").
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
				ID:   e.OrganizationID,
				Name: e.OrganizationName,
			},
			Client: &model.Client{
				ID:   e.ClientID,
				Name: e.ClientName,
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
