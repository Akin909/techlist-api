// Package main_test provides test for the application
package main_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"."
)

var a main.App

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS startups
(
id SERIAL,
name TEXT NOT NULL,
category TEXT NOT NULL,
CONSTRAINT startups_pkey PRIMARY KEY (id)
)`

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize("explorer_test")
	ensureTableExists()
	code := m.Run()

	clearTable()

	os.Exit(code)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ensureTableExists() {
	sql, err := ioutil.ReadFile("./build.sql")
	check(err)
	if _, err := a.DB.Exec(string(sql)); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM startups")
	a.DB.Exec("ALTER SEQUENCE startups_id_seq RESTART WITH 1")
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/startups", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an Empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentStartup(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/startup/fish", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Startup not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Startup not found'. Got '%s'", m["error"])
	}
}

func TestCreateStartup(t *testing.T) {
	clearTable()

	payload := []byte(`{"name:":"test startup", "label": "AI"}`)
	req, _ := http.NewRequest("POST", "/startup", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test startup" {
		t.Errorf("Expected startup name to be 'test startup'. Got '%v'", m["name"])
	}

	if m["label"] != "AI" {
		t.Errorf("Expected startup name to be 'AI'. Got '%v'", m["label"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected startup ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetStartup(t *testing.T) {
	clearTable()
	addStartup(1)

	req, _ := http.NewRequest("GET", "startup", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addStartup(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO startups(name, category) VALUES($1, $2)", "Startup "+strconv.Itoa(i), (i+1.0)*10)
	}
}

func TestUpdateStartup(t *testing.T) {
	clearTable()
	addStartup(1)

	req, _ := http.NewRequest("GET", "/startup/startup", nil)
	response := executeRequest(req)
	var originalStartup map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalStartup)

	payload := []byte(`{"name": "test startup - updated name", "category": "new category"}`)

	req, _ = http.NewRequest("PUT", "/startup/startup", bytes.NewBuffer(payload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalStartup["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalStartup["id"], m["id"])
	}
	if m["name"] == originalStartup["name"] {
		t.Errorf("Expected the name to from '%v' to '%v'. Got '%v'", originalStartup["name"], m["name"], m["name"])
	}
	if m["category"] == originalStartup["category"] {
		t.Errorf("Expected the category to change from '%v' to '%v'. Got '%v'", originalStartup["category"], m["category"], m["category"])
	}
}
