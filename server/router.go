package server

import (
	"database/sql"
	"log"
	"os"
	"path"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/BurntSushi/toml"

	"github.com/shufo/go-graphql-boilerplate/translations"

	"github.com/casbin/casbin/model"

	"github.com/gobuffalo/packr"

	redisadapter "github.com/casbin/redis-adapter"

	"github.com/casbin/casbin"

	"github.com/99designs/gqlgen/handler"
	"github.com/shufo/go-graphql-boilerplate/dataloader"
	"github.com/shufo/go-graphql-boilerplate/graph/generated"
	"github.com/shufo/go-graphql-boilerplate/resolver"
	"github.com/rs/cors"

	"github.com/shufo/go-graphql-boilerplate/auth"
	"github.com/shufo/go-graphql-boilerplate/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/sirupsen/logrus"
)

// Router sets router settings
func (s *Server) Router(db *sql.DB) *chi.Mux {
	/*
	 * Middleware settings
	 */

	// Use JSON logger
	customLogger := logrus.New()
	customLogger.Formatter = &logrus.JSONFormatter{
		// disable, as we set our own
		DisableTimestamp: true,
	}

	// initialize Casbin
	casbin := initCasbin()

	// initialize i18n
	bundle := initI18n()

	// JWT setting
	secret, found := os.LookupEnv("JWT_SECRET")

	if !found {
		log.Fatal("There is no JWT_SECRET variable in environment variables")
	}

	tokenAuth := jwtauth.New("HS256", []byte(secret), nil)

	// middlewares
	s.router.Use(middleware.WithValue("db", db))
	s.router.Use(middleware.WithValue("casbin", casbin))
	s.router.Use(middleware.WithValue("bundle", bundle))
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(jwtauth.Verifier(tokenAuth))
	s.router.Use(translations.Middleware)

	// Logger
	if s.config.Logging {
		s.router.Use(logger.NewStructuredLogger(customLogger))
	}

	s.router.Use(auth.Middleware)
	s.router.Use(middleware.Recoverer)

	allowedOrigins := []string{"*"}

	if env, found := os.LookupEnv("APP_ENV"); found {
		switch env {
		case "production":
			allowedOrigins = []string{"https://www.example.jp"}
		case "development":
			allowedOrigins = []string{"https://dev.www.example.jp", "https://dev.api.example.jp"}
		}
	}

	// CORS setting
	cors := cors.New(cors.Options{
		// Use this to allow specific origin hosts
		AllowedOrigins: allowedOrigins,
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	s.router.Use(cors.Handler)

	// Dataloader for GraphQL
	s.router.Use(dataloader.DataloaderMiddleware)

	/*
	 * Routing settings
	 */

	// GraphQL playground
	s.router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	c := generated.Config{Resolvers: &resolver.Resolver{}, Directives: resolver.NewDirectives()}

	// GraphQL endpoint
	s.router.Handle("/query", handler.GraphQL(generated.NewExecutableSchema(c)))

	return s.router
}

func initCasbin() *casbin.CachedEnforcer {
	// redis adapter
	host := os.Getenv("REDIS_HOST")
	a := redisadapter.NewAdapter("tcp", host+":6379")

	// load config from packr
	box := packr.NewBox("../configs")
	modelText, err := box.FindString("casbin_rbac.conf")

	if err != nil {
		log.Fatal("casbin model not found")
	}

	model := model.Model{}
	model.LoadModelFromText(modelText)

	// use cached enforcer
	e := casbin.NewCachedEnforcer(model, a)

	// Load the policy from DB.
	e.LoadPolicy()

	return e
}

func initI18n() *i18n.Bundle {
	// Init i18n package
	bundle := &i18n.Bundle{DefaultLanguage: language.English}
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load translation files from packr
	box := packr.NewBox("../translations")
	messageFiles := box.List()

	for _, file := range messageFiles {
		if path.Ext(file) == ".toml" {
			bytes, _ := box.Find(file)
			bundle.MustParseMessageFileBytes(bytes, file)
		}
	}

	return bundle
}
