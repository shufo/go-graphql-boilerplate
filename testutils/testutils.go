package testutils

import (
	"database/sql"
	"log"
	"os"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures"
	"github.com/shufo/go-graphql-boilerplate/server"
	"github.com/shufo/go-graphql-boilerplate/utils"
)

func PrepareDB() *sql.DB {
	// get environment variables
	hostname := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")

	dsn := username + ":" + password + "@tcp(" + hostname + ":" + port + ")/" + "?charset=utf8mb4&parseTime=True&loc=Local"

	// connect to database
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("can't conect to MySQL")
	}

	// create db if not exists
	createSQL := "CREATE DATABASE IF NOT EXISTS " + database + "_test " + "CHARACTER SET utf8mb4"

	if res, err := db.Exec(createSQL); err != nil {
		log.Fatal(res)
	}

	db.Exec("use " + database + "_test")

	// migrate db
	MigrateDB(db)

	dsn = username + ":" + password + "@tcp(" + hostname + ":" + port + ")/" + database + "_test?charset=utf8mb4&parseTime=True&loc=Local"

	var alreadyRegisted bool

	// avoid driver registry duplication
	for _, driver := range sql.Drivers() {
		if driver == "txdb" {
			alreadyRegisted = true
		}
	}

	if !alreadyRegisted {
		txdb.Register("txdb", "mysql", dsn)
	}

	// open db with transactional connection
	db, _ = sql.Open("txdb", dsn)

	/* Uncomment this if you want to display queries on test
	boil.DebugMode = true
	*/

	return db
}

// MigrateDB migrates primary DB instance with passed db connection
func MigrateDB(db *sql.DB) {
	utils.MigrateDB(db)
}

func PrepareRouter(db *sql.DB) *chi.Mux {
	// prepare router for testing
	c := server.Config{Logging: false}
	s := server.NewServer(c)
	r := s.Router(db)

	return r
}

func PopulateRecords(db *sql.DB) {
	err := testfixtures.GenerateFixtures(db, &testfixtures.MySQL{}, "../fixtures")
	if err != nil {
		log.Fatalf("Error generating fixtures: %v", err)
	}
}
