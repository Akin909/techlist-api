// Package main provides an entry point to the app
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// App references the router and the database the app uses
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Options represent the server config options
type Options struct {
	path string
	port string
}

const (
	dbUser     = "A_nonymous"
	dbPassword = "postgres"
)

// EnsureTableExists creates the startup table if it doesn't already exist
func EnsureTableExists(db *sql.DB) {
	sql, err := ioutil.ReadFile("./build.sql")
	check(err)

	if _, err := db.Exec(string(sql)); err != nil {
		log.Fatal(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// getStartup gets the variables from the request respondWithError if any err otherwise it responds with the result
func (a *App) getStartup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid startup ID")
		return
	}

	s := startup{ID: id}
	if err := s.getStartup(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Startup not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, s)
}

// respondWithJSON Marshals the data to be sent sets the header and writes the response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError returns a JSON object with an error key and error message value
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func (a *App) getStartups(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}

	if start < 0 {
		start = 0
	}

	startups, err := getStartups(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, startups)
}

func (a *App) createStartup(w http.ResponseWriter, r *http.Request) {
	var s startup
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := s.createStartup(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, s)
}

// Initialize function starts the application
func (a *App) Initialize(name string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, name)
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.InitializeRoutes()

}

func (a *App) updateStartup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	var s startup
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	defer r.Body.Close()
	s.ID = id

	if err := s.updateStartup(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, s)
}

// InitializeRoutes creates handlers for the different application routes
func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/startups", a.getStartups).Methods("GET")
	a.Router.HandleFunc("/startup", a.createStartup).Methods("POST")
	a.Router.HandleFunc("/startup/{id:[0-9]+}", a.getStartup).Methods("GET")
	a.Router.HandleFunc("/startup/{id:[0-9]+}", a.updateStartup).Methods("PUT")
}

// Log wraps each handler function in a logger function which prints
// the remote address, method and url
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// Run function starts the app on a given port
func (a *App) Run(addr string) {
	op := &Options{port: ":8001"}
	fmt.Printf("Looking out for new startups on %s", op.port)
	err := http.ListenAndServe(op.port, Log(a.Router))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
