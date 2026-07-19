package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"finance-app-backend/internal/config"
	"finance-app-backend/internal/database"
	"finance-app-backend/internal/database/migrations"
	"finance-app-backend/internal/database/seeders"
)

func main() {
	config.Load()
	database.Connect()

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "migrate":
		if err := migrations.Migrate(database.DB); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "migrate:rollback":
		if err := migrations.Rollback(database.DB); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
	case "migrate:fresh":
		forceFlag := false
		for _, arg := range os.Args[2:] {
			if arg == "--force" {
				forceFlag = true
				break
			}
		}

		if !forceFlag && config.AppConfig.AppEnv == "production" {
			log.Fatal("Refusing to run migrate:fresh in production without --force flag.")
		}

		if err := migrations.Fresh(database.DB); err != nil {
			log.Fatalf("Fresh migration failed: %v", err)
		}
	case "migrate:status":
		if err := migrations.Status(database.DB); err != nil {
			log.Fatalf("Migration status failed: %v", err)
		}
	case "seed":
		seedFlags := flag.NewFlagSet("seed", flag.ExitOnError)
		onlySeeder := seedFlags.String("only", "", "Specify a single seeder to run")
		seedFlags.Parse(os.Args[2:])

		// Fallback parse if passed like --only=UserSeeder without flag parsing
		if *onlySeeder == "" {
			for _, arg := range os.Args[2:] {
				if strings.HasPrefix(arg, "--only=") {
					*onlySeeder = strings.TrimPrefix(arg, "--only=")
				}
			}
		}

		if err := seeders.RunAll(database.DB, *onlySeeder); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Personal Finance CLI Tool")
	fmt.Println("Usage: go run ./cli [command] [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  migrate                  - Run pending database migrations")
	fmt.Println("  migrate:rollback         - Rollback last batch of migrations")
	fmt.Println("  migrate:fresh [--force]  - Drop all tables and run migrations")
	fmt.Println("  migrate:status           - Display migration status")
	fmt.Println("  seed [--only=SeederName] - Run registered database seeders")
}
