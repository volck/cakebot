package main

import (
	"cakebot/internal/data"
	"cakebot/internal/notifier"
	"context"
	"database/sql"
	b64 "encoding/base64"
	"flag"
	"fmt"
	_ "github.com/godror/godror"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"net/http"
	"os"
	"time"
)

// Declare a string containing the application version number. Later in the book we'll
// generate this automatically at build time, but for now we'll just store the version
// number as a hard-coded global constant.
const version = "1.0.0"

// Define a config struct to hold all the configuration settings for our application.
// For now, the only configuration settings will be the network port that we want the
// server to listen on, and the name of the current operating environment for the
// application (development, staging, production, etc.). We will read in these
// configuration settings from command-line flags when the application starts.
type config struct {
	port int
	env  string
	db   struct {
		jdbc_url string
	}
	keycloak           oauth2.Config
	graphApis          clientcredentials.Config
	graphApiCalendarID string
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config   config
	logger   *log.Logger
	models   data.Models
	notifier notifier.Notifier
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.keycloak.ClientID, "clientId", os.Getenv("KEYCLOAK_CLIENTID"), "Keycloak client id")
	flag.StringVar(&cfg.keycloak.ClientSecret, "clientSecret", os.Getenv("KEYCLOAK_CLIENTSECRET"), "Keycloak client id")
	flag.StringVar(&cfg.keycloak.Endpoint.TokenURL, "tokenUrl", os.Getenv("KEYCLOAK_TOKENURL"), "Keycloak token url")
	flag.StringVar(&cfg.keycloak.Endpoint.AuthURL, "AuthURL", os.Getenv("KEYCLOAK_AUTHURL"), "Keycloak token url")
	flag.StringVar(&cfg.graphApis.ClientID, "graphClientId", os.Getenv("GRAPH_CLIENTID"), "GraphApi client id")
	flag.StringVar(&cfg.graphApis.ClientSecret, "graphClientSecret", os.Getenv("GRAPH_CLIENTSECRET"), "GraphApi client id")
	flag.StringVar(&cfg.graphApis.TokenURL, "graphTokenUrl", os.Getenv("GRAPH_TOKENURL"), "GraphApi token URL")
	flag.StringVar(&cfg.graphApiCalendarID, "graphCalendarID", os.Getenv("GRAPH_CALENDARID"), "GraphApi calendar ID")

	flag.StringVar(&cfg.db.jdbc_url, "jdbc", os.Getenv("JDBC_URL"), "Oracle JDBC URL")
	flag.Parse()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	keycloakConfig := createKeyCloakClient(cfg)

	keycloakHTTPClient := data.KeycloakClient{Config: keycloakConfig}
	logger.Printf("database connection pool established")
	// Use the data.NewModels() function to initialize a Models struct, passing in the
	// connection pool as a parameter.

	app := &application{
		config:   cfg,
		logger:   logger,
		models:   data.NewModels(db, keycloakHTTPClient),
		notifier: notifier.New(cfg.graphApis),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go app.notifier.Notify(app.models.Cake, app.models.Cake)
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.

	jdbcDecoded, err := b64.StdEncoding.DecodeString(cfg.db.jdbc_url)
	if err != nil {
		fmt.Println("err", err)
	}
	db, err := sql.Open("godror", string(jdbcDecoded))
	if err != nil {
		return nil, err
	}
	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	// Return the sql.DB connection pool.
	return db, nil
}

// returns a valid oauth2 config for interfacing with keycloak
func createKeyCloakClient(cfg config) oauth2.Config {
	keycloakConfig := oauth2.Config{
		ClientID:     cfg.keycloak.ClientID,
		ClientSecret: cfg.keycloak.ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: cfg.keycloak.Endpoint.TokenURL,
			AuthURL:  cfg.keycloak.Endpoint.AuthURL,
		},
	}
	return keycloakConfig

}
