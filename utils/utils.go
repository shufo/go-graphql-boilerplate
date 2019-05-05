package utils

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/packr"
	migrate "github.com/rubenv/sql-migrate"
)

// MigrateDB migrates primary DB instance with passed db connection
func MigrateDB(db *sql.DB) {
	migrationSource := &migrate.PackrMigrationSource{
		Box: packr.NewBox("../migrations"),
	}

	fmt.Printf("[migration] Begin migration...\n\n")

	i, err := migrate.Exec(db, "mysql", migrationSource, migrate.Up)

	if err != nil {
		log.Fatal(err)
	}

	migrations, _ := migrationSource.FindMigrations()

	if len(migrations) > 0 {
		fmt.Printf("[migration] Last migration id: %s\n\n", migrations[len(migrations)-1].Id)
	}

	if i > 0 {
		fmt.Printf("[migration] %d migrations executed\n\n", i)
	} else {
		fmt.Printf("[migration] There is no new migration\n\n")
	}

	fmt.Printf("[migration] Migration finished\n\n")
}

func RandomUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		log.Fatal(err)
	}

	return uuid.New().String()
}
