package repository

import "github.com/take73/invoice-api-example/internal/domain/model"

type InvoiceRepository interface {
	Create(invoice model.Invoice) (*model.Invoice, error)
}
