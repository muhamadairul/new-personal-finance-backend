package seeders

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// RunAll executes all registered seeders or a specific seeder if specified
func RunAll(db *gorm.DB, onlySeeder string) error {
	executed := 0
	for _, s := range registry {
		if onlySeeder != "" && s.Name() != onlySeeder {
			continue
		}

		log.Printf("Seeding: %s", s.Name())
		if err := s.Run(db); err != nil {
			return fmt.Errorf("seeder %s failed: %w", s.Name(), err)
		}
		log.Printf("Seeded:  %s", s.Name())
		executed++
	}

	if executed == 0 {
		if onlySeeder != "" {
			return fmt.Errorf("seeder '%s' not found in registry", onlySeeder)
		}
		log.Println("No seeders registered.")
	}

	return nil
}
