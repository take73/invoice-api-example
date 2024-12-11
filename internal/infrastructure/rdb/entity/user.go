package entity

import "time"

// User ORMのEntity
type User struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;column:user_id"`
	OrganizationID uint      `gorm:"column:organization_id;not null"`
	Name           string    `gorm:"column:name;not null"`
	Email          string    `gorm:"column:email;not null;unique"`
	Password       string    `gorm:"column:password;not null"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Associations
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name used by GORM.
func (User) TableName() string {
	return "user"
}
