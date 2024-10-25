package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type todoStorage interface {
	ReadTodos() ([]Todo, error)
	CreateTodo(subject string, details string, priority int, dueToDate string) error
	DeleteTodo(id string) error
	// TODO: Add function to mark TODO as DONE
}

type StorageImpl struct {
	connection *sql.DB
}

func NewStorage(dbType, dbName string) (StorageImpl, error) {
	connection, err := sql.Open(dbType, dbName)
	if err != nil {
		return StorageImpl{}, err
	}
	log.Printf("Connected to database %v", connection)
	return StorageImpl{
		connection: connection,
	}, nil
}

func (s StorageImpl) ReadTodos() ([]Todo, error) {
	rows, err := s.connection.Query("SELECT id, subject, details, priority, due_to_date, resolved FROM todos ORDER BY id")
	if err != nil {
		return []Todo{}, err
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var id int
		var subject string
		var details string
		var priority int
		var dueToDate string
		var resolved bool

		if err := rows.Scan(&id, &subject, &details, &priority, &dueToDate, &resolved); err != nil {
			return todos, err
		}
		todos = append(todos, Todo{
			ID:        id,
			Subject:   subject,
			Details:   details,
			Priority:  priority,
			DueToDate: dueToDate,
			Resolved:  resolved,
		})
	}

	return todos, nil
}

func (s StorageImpl) CreateTodo(subject string, details string, priority int, dueToDate string) error {
	statement, err := s.connection.Prepare("INSERT INTO todos(subject, details, priority, due_to_date, resolved) VALUES (?, ?, ?, ?, 0)")
	_, err = statement.Exec(subject, details, priority, dueToDate)
	return err
}

func (s StorageImpl) DeleteTodo(id string) error {
	statement, err := s.connection.Prepare("DELETE FROM todos WHERE id=?")
	_, err = statement.Exec(id)
	return err
}
