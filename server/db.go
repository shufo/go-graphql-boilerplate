package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/volatiletech/sqlboiler/boil"
)

func (s *Server) OpenDBConnection() *sql.DB {
	//Check if environment variable exists
	envs := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USERNAME",
		"DB_PASSWORD",
		"DB_DATABASE",
	}

	for _, env := range envs {
		if _, found := os.LookupEnv(env); !found {
			log.Fatalf("Required environment variable %s not found.", env)
		}
	}

	hostname := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")

	dsn := username + ":" + password + "@tcp(" + hostname + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"

	var db *sql.DB
	var err error

	// try connecting to db for several times
	for connectionCount := 1; connectionCount < 5; connectionCount++ {
		db, err = sql.Open("mysql", dsn)

		if err := db.Ping(); err != nil {
			fmt.Printf("can't connect to database: %v. reconnecting...\n\n", err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	if err != nil {
		fmt.Printf("%v", err)
		panic("failed to connect database")
	}

	if env, found := os.LookupEnv("APP_ENV"); !found || env != "production" {
		boil.DebugMode = true
	}

	return db
}
