package migrations

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// Migration represents a single migration unit
type Migration struct {
	Name string
	Up   func(db *gorm.DB) error
	Down func(db *gorm.DB) error
}

// MigrationRecord represents the migrations tracking table
type MigrationRecord struct {
	ID        uint      `gorm:"primaryKey"`
	Migration string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Batch     int       `gorm:"not null"`
	AppliedAt time.Time `gorm:"autoCreateTime"`
}

// TableName explicitly names the tracking table
func (MigrationRecord) TableName() string {
	return "migrations"
}

var registry []Migration

// Register registers a new migration file
func Register(m Migration) {
	registry = append(registry, m)
}

// EnsureTable ensures the migrations tracking table exists
func EnsureTable(db *gorm.DB) error {
	return db.AutoMigrate(&MigrationRecord{})
}

// Migrate runs all pending migrations
func Migrate(db *gorm.DB) error {
	if err := EnsureTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	var applied []MigrationRecord
	if err := db.Find(&applied).Error; err != nil {
		return err
	}

	appliedMap := make(map[string]bool)
	maxBatch := 0
	for _, rec := range applied {
		appliedMap[rec.Migration] = true
		if rec.Batch > maxBatch {
			maxBatch = rec.Batch
		}
	}

	currentBatch := maxBatch + 1
	executedCount := 0

	for _, m := range registry {
		if appliedMap[m.Name] {
			continue
		}

		log.Printf("Migrating: %s", m.Name)
		if err := m.Up(db); err != nil {
			return fmt.Errorf("failed migration %s: %w", m.Name, err)
		}

		rec := MigrationRecord{
			Migration: m.Name,
			Batch:     currentBatch,
			AppliedAt: time.Now(),
		}
		if err := db.Create(&rec).Error; err != nil {
			return fmt.Errorf("failed to record migration %s: %w", m.Name, err)
		}

		log.Printf("Migrated:  %s", m.Name)
		executedCount++
	}

	if executedCount == 0 {
		log.Println("Nothing to migrate.")
	} else {
		log.Printf("Successfully ran %d migration(s) (Batch %d).", executedCount, currentBatch)
	}

	return nil
}

// Rollback rolls back the last batch of migrations
func Rollback(db *gorm.DB) error {
	if err := EnsureTable(db); err != nil {
		return err
	}

	var maxBatch int
	row := db.Model(&MigrationRecord{}).Select("COALESCE(MAX(batch), 0)").Row()
	if err := row.Scan(&maxBatch); err != nil || maxBatch == 0 {
		log.Println("Nothing to rollback.")
		return nil
	}

	var lastBatchRecords []MigrationRecord
	if err := db.Where("batch = ?", maxBatch).Order("id DESC").Find(&lastBatchRecords).Error; err != nil {
		return err
	}

	registryMap := make(map[string]Migration)
	for _, m := range registry {
		registryMap[m.Name] = m
	}

	for _, rec := range lastBatchRecords {
		m, ok := registryMap[rec.Migration]
		if !ok {
			log.Printf("Warning: Migration %s in DB not found in registry", rec.Migration)
			continue
		}

		log.Printf("Rolling back: %s", m.Name)
		if err := m.Down(db); err != nil {
			return fmt.Errorf("failed rollback %s: %w", m.Name, err)
		}

		if err := db.Delete(&rec).Error; err != nil {
			return fmt.Errorf("failed to delete migration record %s: %w", m.Name, err)
		}

		log.Printf("Rolled back:  %s", m.Name)
	}

	log.Printf("Rolled back Batch %d.", maxBatch)
	return nil
}

// Fresh drops all database tables and re-runs all migrations
func Fresh(db *gorm.DB) error {
	log.Println("Dropping all tables...")

	// Disable FK checks if MySQL / Postgres
	db.Exec("SET FOREIGN_KEY_CHECKS = 0;") // MySQL
	defer db.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	tables, err := db.Migrator().GetTables()
	if err == nil {
		for _, t := range tables {
			if err := db.Migrator().DropTable(t); err != nil {
				log.Printf("Failed to drop table %s: %v", t, err)
			}
		}
	}

	log.Println("All tables dropped. Running fresh migrations...")
	return Migrate(db)
}

// Status displays the status of all registered migrations
func Status(db *gorm.DB) error {
	if err := EnsureTable(db); err != nil {
		return err
	}

	var applied []MigrationRecord
	if err := db.Find(&applied).Error; err != nil {
		return err
	}

	appliedMap := make(map[string]MigrationRecord)
	for _, rec := range applied {
		appliedMap[rec.Migration] = rec
	}

	fmt.Println("\n+------------------------------------------------------+---------+-------+")
	fmt.Println("| Migration File Name                                  | Status  | Batch |")
	fmt.Println("+------------------------------------------------------+---------+-------+")

	for _, m := range registry {
		if rec, ok := appliedMap[m.Name]; ok {
			fmt.Printf("| %-52s | APPLIED | %-5d |\n", m.Name, rec.Batch)
		} else {
			fmt.Printf("| %-52s | PENDING | -     |\n", m.Name)
		}
	}
	fmt.Println("+------------------------------------------------------+---------+-------+\n")

	return nil
}
