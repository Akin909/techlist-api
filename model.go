package main

import "database/sql"

type startup struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Catergory string `json:"price"`
}

func (s *startup) getStartup(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM startup WHERE id=$1", s.ID).Scan(&s.Name, &s.Catergory)
}

func (s *startup) updateStartup(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE startups SET name=$1, category=$2 WHERE id=$3", s.Name, s.Catergory, s.ID)
	return err
}

func (s *startup) deleteStartup(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM startups WHERE id=$3", s.Name, s.Catergory, s.ID)
	return err
}

func (s *startup) createStartup(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO startups (name, category) VALUES($1, $2) RETURNING id", s.Name, s.Catergory).Scan(&s.ID)
	if err != nil {
		return err
	}
	return nil
}

func getStartups(db *sql.DB, start, count int) ([]startup, error) {
	rows, err := db.Query("SELECT id, name, category FROM startups LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	startups := []startup{}

	for rows.Next() {
		var s startup
		if err := rows.Scan(&s.ID, &s.Name, &s.Catergory); err != nil {
			return nil, err
		}
		startups = append(startups, s)
	}
	return startups, nil
}
