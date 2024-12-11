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

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) repository.Organization {
	return &OrganizationRepository{db: db}
}

// GetByID retrieves a client by its ID.
func (r *OrganizationRepository) GetByID(id uint) (*model.Organization, error) {
	var entity entity.Organization
	if err := r.db.Where("organization_id = ?", id).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commonErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to retrieve organization with ID %d: %w", id, err)
	}

	organization := &model.Organization{
		ID:             entity.ID,
		Name:           entity.Name,
		Representative: entity.RepresentativeName,
		PhoneNumber:    entity.PhoneNumber,
		PostalCode:     entity.PostalCode,
		Address:        entity.Address,
	}

	return organization, nil
}
