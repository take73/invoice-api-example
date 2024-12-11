package rdb

import (
	"errors"
	"fmt"

	"github.com/take73/invoice-api-example/internal/domain/model"
	"github.com/take73/invoice-api-example/internal/domain/repository"
	"github.com/take73/invoice-api-example/internal/infrastructure/rdb/entity"
	commonErrors "github.com/take73/invoice-api-example/internal/shared/errors"
	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) repository.Client {
	return &ClientRepository{db: db}
}

// GetByID retrieves a client by its ID.
func (r *ClientRepository) GetByID(id uint) (*model.Client, error) {
	var entity entity.Client
	if err := r.db.Preload("Organization").Where("client_id = ?", id).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commonErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to retrieve client with ID %d: %w", id, err)
	}

	client := &model.Client{
		ID:             entity.ID,
		OrganizationID: entity.OrganizationID,
		Name:           entity.Name,
		Representative: entity.RepresentativeName,
		PhoneNumber:    entity.PhoneNumber,
		PostalCode:     entity.PostalCode,
		Address:        entity.Address,
	}

	return client, nil
}
