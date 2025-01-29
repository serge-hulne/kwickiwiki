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
