package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type todoStorage interface {
	ReadTodos() ([]Todo, error)
	CreateTodo(name string, content string) error
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
	rows, err := s.connection.Query("SELECT id, name, content, resolved FROM todos ORDER BY id")
	if err != nil {
		return []Todo{}, err
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var id int
		var name string
		var content string
		var resolved bool

		if err := rows.Scan(&id, &name, &content, &resolved); err != nil {
			return todos, err
		}
		todos = append(todos, Todo{
			ID:       id,
			Name:     name,
			Content:  content,
			Resolved: resolved,
		})
	}

	return todos, nil
}

func (s StorageImpl) CreateTodo(name string, surname string) error {
	statement, err := s.connection.Prepare("INSERT INTO todos(name, content, false) VALUES (?, ?, ?)")
	_, err = statement.Exec(name, surname)
	return err
}

func (s StorageImpl) DeleteTodo(id string) error {
	statement, err := s.connection.Prepare("DELETE FROM todos WHERE id=?")
	_, err = statement.Exec(id)
	return err
}
