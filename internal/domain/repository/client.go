package repository

import "github.com/take73/invoice-api-example/internal/domain/model"

type Client interface {
	GetByID(id uint) (*model.Client, error)
}
