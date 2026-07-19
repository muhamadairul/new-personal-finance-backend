package seeders

import (
	"gorm.io/gorm"
)

// Seeder interface defines how seeders execute
type Seeder interface {
	Name() string
	Run(db *gorm.DB) error
}

var registry []Seeder

// Register registers a seeder into global registry
func Register(s Seeder) {
	registry = append(registry, s)
}
