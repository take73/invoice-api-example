package entity

import "time"

// Client ORMのEntity
type Client struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement;column:client_id"`
	OrganizationID     uint      `gorm:"column:organization_id;not null"`
	Name               string    `gorm:"column:name;not null"`
	RepresentativeName string    `gorm:"column:representative_name;not null"`
	PhoneNumber        string    `gorm:"column:phone_number"`
	PostalCode         string    `gorm:"column:postal_code"`
	Address            string    `gorm:"column:address"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Associations
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name used by GORM.
func (Client) TableName() string {
	return "client"
}
