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

func (r *OrganizationRepository) GetByUserID(userID uint) (*model.Organization, error) {
	// joinしたクエリ結果を格納する構造体
	type result struct {
		OrganizationID     uint
		OrganizationName   string
		RepresentativeName string
		PhoneNumber        string
		PostalCode         string
		Address            string
	}

	var res result

	if err := r.db.Table("user").
		Select("organization.organization_id, organization.name AS organization_name, "+
			"organization.representative_name, organization.phone_number, "+
			"organization.postal_code, organization.address").
		Joins("JOIN organization ON user.organization_id = organization.organization_id").
		Where("user.user_id = ?", userID).
		Scan(&res).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commonErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to retrieve organization for user ID %d: %w", userID, err)
	}

	organization := &model.Organization{
		ID:             res.OrganizationID,
		Name:           res.OrganizationName,
		Representative: res.RepresentativeName,
		PhoneNumber:    res.PhoneNumber,
		PostalCode:     res.PostalCode,
		Address:        res.Address,
	}

	return organization, nil
}
