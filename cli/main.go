package main

import (
	"fmt"
	"os"

	"finance-app-backend/internal/config"
)

func main() {
	config.Load()

	if len(os.Args) < 2 {
		fmt.Println("Usage: cli [command]")
		fmt.Println("Commands:")
		fmt.Println("  migrate          - Run pending database migrations")
		fmt.Println("  migrate:rollback - Rollback last batch of migrations")
		fmt.Println("  migrate:fresh    - Drop all tables and run migrations")
		fmt.Println("  migrate:status   - Display migration status")
		fmt.Println("  seed             - Seed data using registered seeders")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "migrate":
		fmt.Println("Running migrations...")
		// TODO: Implement migration execution
	case "migrate:rollback":
		fmt.Println("Rolling back last migration batch...")
		// TODO: Implement migration rollback
	case "migrate:fresh":
		fmt.Println("Fresh migration: dropping and recreating all tables...")
		// TODO: Implement fresh migration
	case "migrate:status":
		fmt.Println("Checking migration status...")
		// TODO: Implement status check
	case "seed":
		fmt.Println("Seeding database...")
		// TODO: Implement seed execution
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
