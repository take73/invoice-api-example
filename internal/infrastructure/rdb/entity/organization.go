package entity

import "time"

// Organization ORMのEntity
type Organization struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement;column:organization_id"`
	Name               string    `gorm:"column:name;not null"`
	RepresentativeName string    `gorm:"column:representative_name;not null"`
	PhoneNumber        string    `gorm:"column:phone_number"`
	PostalCode         string    `gorm:"column:postal_code"`
	Address            string    `gorm:"column:address"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Associations
	Users   []User   `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Clients []Client `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name used by GORM.
func (Organization) TableName() string {
	return "organization"
}
