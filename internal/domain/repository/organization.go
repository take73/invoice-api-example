package repository

import "github.com/take73/invoice-api-example/internal/domain/model"

type Organization interface {
	GetByID(id uint) (*model.Organization, error)
}
