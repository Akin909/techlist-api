// Package main_test provides test for the application
package main_test

import (
	"log"
	"os"
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
	a.Initialize()
	ensureTableExists()
	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM startups")
	a.DB.Exec("ALTER SEQUENCE startups_id_seq RESTART WITH 1")
}
