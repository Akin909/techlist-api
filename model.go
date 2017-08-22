// Package main provides access to the DB information
package main

import (
	"database/sql"
	"errors"
)

type product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

func (p *product) getStartup(db *sql.DB) error {
	return errors.New("Not Implemented")
}

func (p *product) updateStartup(db *sql.DB) error {
	return errors.New("Not Implemented")
}

func (p *product) deleteStartup(db *sql.DB) error {
	return errors.New("Not Implemented")
}

func (p *product) createStartup(db *sql.DB) error {
	return errors.New("Not Implemented")
}

func getStartups(db *sql.DB, start, count int) ([]product, error) {
	return nil, errors.New("Not Implemented")
}
