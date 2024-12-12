package repository

import (
	"time"

	"github.com/take73/invoice-api-example/internal/domain/model"
)

type Invoice interface {
	Create(invoice *model.Invoice) (*model.Invoice, error)
	FindByDueDateRange(startDate, endDate time.Time) ([]*model.Invoice, error)
}
