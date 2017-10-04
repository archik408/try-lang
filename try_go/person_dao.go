package main

import (
	"database/sql"
)

type person struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Data string `json:"data"`
}

func (p *person) getPerson(db *sql.DB) error {
	return db.QueryRow("SELECT name, data FROM person WHERE id=$1",
		p.ID).Scan(&p.Name, &p.Data)
}

func (p *person) updatePerson(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE person SET name=$1, data=$2 WHERE id=$3",
			p.Name, p.Data, p.ID)

	return err
}

func (p *person) deletePerson(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM person WHERE id=$1", p.ID)

	return err
}

func (p *person) createPerson(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO person(name, data) VALUES($1, $2) RETURNING id",
		p.Name, p.Data).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}


func getPersons(db *sql.DB, start, count int) ([]person, error) {
	rows, err := db.Query(
		"SELECT id, name, data FROM person LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	persons := []person{}

	for rows.Next() {
		var p person
		if err := rows.Scan(&p.ID, &p.Name, &p.Data); err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}

	return persons, nil
}
