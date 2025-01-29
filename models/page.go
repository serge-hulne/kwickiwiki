package models

import (
	"gorm.io/datatypes" // Required for JSON fields in GORM
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model
	Title    string `gorm:"unique"`
	Content  string
	Metadata datatypes.JSON `gorm:"type:json"` // JSON field for all metadata
}

type User struct {
	gorm.Model
	Name     string `gorm:"uniqueIndex"`
	Email    string `gorm:"uniqueIndex"`
	Password string
	Roles    []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
	gorm.Model
	Name string `gorm:"uniqueIndex"`
}

type UserRole struct {
	gorm.Model
	UserID uint
	RoleID uint
}
