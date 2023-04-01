package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/snirkop89/go-movies/internal/repository"
	"github.com/snirkop89/go-movies/internal/repository/dbrepo"
	"github.com/snirkop89/simplelogger"
)

const port = 8080

type application struct {
	logger       *simplelogger.Logger
	DSN          string
	Domain       string
	DB           repository.DatabaseRepo
	auth         auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	APIKey       string
}

func main() {
	// set application config
	var app application

	// Config options from command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5437 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.APIKey, "api-key", "c628aab7009d82e6a615f654e8fbda33", "The movieDB API key")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.Parse()

	// Initialize logger
	app.logger = simplelogger.New(simplelogger.FormatHuman, simplelogger.LevelInfo)

	// connect to database
	conn, err := app.connectToDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	// TODO - create simple logger package

	// TODO Replace the wasterful vars
	app.auth = auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "Host-refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	// start a webserver
	app.logger.Infof("Starting application on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
