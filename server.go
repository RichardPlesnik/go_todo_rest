package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

type Server interface {
	Serve(port uint)
}

type ServerImpl struct {
	todoStorage todoStorage
}

func NewServer(todoStorage todoStorage) Server {
	return ServerImpl{
		todoStorage: todoStorage,
	}
}

func (s ServerImpl) indexPageHandler(writer http.ResponseWriter, r *http.Request) {
	http.ServeFile(writer, r, "index.html")
}

func (s ServerImpl) listUsersPageHandler(writer http.ResponseWriter, r *http.Request) {
	// construct template on the fly - allow us to change the template
	// while the service is running
	const templateFilename = "todos.html"
	log.Printf("Constructing template from file %s", templateFilename)
	// new template
	tmpl, err := template.ParseFiles(templateFilename)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("Template can't be constructed: %v", err)
		return
	}

	users, err := s.todoStorage.ReadTodos()
	if err != nil {
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusInternalServerError)
		if err != nil {
			log.Printf("Unable to retrieve list of todos: %v", err)
		}
		_, err := writer.Write([]byte("Unable to retrieve list of todos"))
		if err != nil {
			log.Printf("Unable to retrieve list of todos: %v", err)
		}
		return
	}
	log.Printf("Application template for %d data records", len(users))

	// apply template
	err = tmpl.Execute(writer, users)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}

func (s ServerImpl) createTodoPageHandler(writer http.ResponseWriter, r *http.Request) {
	http.ServeFile(writer, r, "create_todo.html")
}

func (s ServerImpl) createTodoHandler(writer http.ResponseWriter, r *http.Request) {
	subject := r.FormValue("subject")
	details := r.FormValue("details")
	priority := r.FormValue("priority")
	priority_value, err := strconv.Atoi(priority)
	if err != nil {
		priority_value = -1
		log.Println("Setting default priority -1")

	}
	dueToDate := r.FormValue("due_to_date")

	s.todoStorage.CreateTodo(subject, details, priority_value, dueToDate)
	log.Println("Creating new todo", subject, details, priority_value, dueToDate)
	http.ServeFile(writer, r, "index.html")
}

func (s ServerImpl) usersAPIHandler(writer http.ResponseWriter, r *http.Request) {
	users, err := s.todoStorage.ReadTodos()
	if err != nil {
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusInternalServerError)
		if err != nil {
			log.Printf("Unable to retrieve list of todos: %v", err)
		}
		_, err := writer.Write([]byte("Unable to retrieve list of todos"))
		if err != nil {
			log.Printf("Unable to retrieve list of todos: %v", err)
		}
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(users)
}

// func (s ServerImpl) deleteUserAPIHandler(writer http.ResponseWriter, r *http.Request) {
// 	userID := r.PathValue("id")
// 	log.Println("Going to delete user with ID", userID)
// 	s.userStorage.DeleteUser(userID)
// }

// startServer starts HTTP server that provides all static and dynamic data
func (s ServerImpl) Serve(port uint) {
	log.Printf("Starting server on port %d", port)
	// HTTP pages
	http.HandleFunc("/", s.indexPageHandler)
	http.HandleFunc("/list-todos", s.listUsersPageHandler)
	http.HandleFunc("/create-todo-form", s.createTodoPageHandler)
	http.HandleFunc("/create-todo", s.createTodoHandler)

	// REST API endpoints
	http.HandleFunc("/todos", s.usersAPIHandler)
	// http.HandleFunc("DELETE /user/{id}", s.deleteUserAPIHandler)

	// start the server
	http.ListenAndServe(":8080", nil)
}
