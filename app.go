// Package main provides an entry point to the app
package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
)

// App references the router and the database the app uses
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

const (
	dbUser     = "A_nonymous"
	dbPassword = "postgres"
	dbName     = "test"
)

// Initialize function starts the application
func (a *App) Initialize() {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
}

// Run function starts the app on a given port
func (a *App) Run(addr string) {

}
