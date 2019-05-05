package main

import (
	"log"
	"net/http"
	"os"

	"github.com/shufo/go-graphql-boilerplate/server"
	"github.com/shufo/go-graphql-boilerplate/utils"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = defaultPort
	}

	c := server.Config{Logging: true}
	s := server.NewServer(c)
	db := s.OpenDBConnection()
	utils.MigrateDB(db)
	r := s.Router(db)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
