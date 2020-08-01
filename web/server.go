package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/i-hate-nicknames/tasker/db"
)

type ctxKey int

const (
	keyDb ctxKey = iota
)

func StartServer() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	var db db.Database = db.MakeMemoryDb()
	r.Use(middleware.WithValue(keyDb, db))

	r.Route("/projects", func(r chi.Router) {
		r.Get("/", listProjects)
		r.Post("/", createProject)

		r.Route("/{projectID}", func(r chi.Router) {
			r.Get("/", getProject)
			r.Get("/columns", getColumns)
			r.Put("/", updateProject)
			r.Delete("/", deleteProject)
		})
	})

	r.Route("/columns", func(r chi.Router) {
		r.Post("/", createColumn)
		r.Get("/swap/{orig}/{target}", swapColumns)

		r.Route("/{columnID}", func(r chi.Router) {
			r.Put("/", updateColumn)
			r.Delete("/", deleteColumn)
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("chiii"))
	})
	http.ListenAndServe(":3000", r)

}

func listProjects(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(keyDb).(db.Database)
	projects, err := db.GetProjects()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(projects)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func getProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	ID, err := strconv.Atoi(projectID)
	if err != nil {
		http.Error(w, "project id must be an int", http.StatusBadRequest)
		return
	}
	db := r.Context().Value(keyDb).(db.Database)
	project, err := db.GetProject(ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("project %d not found", ID), http.StatusNotFound)
		return
	}
	response, err := json.Marshal(project)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func createProject(w http.ResponseWriter, r *http.Request) {

}

func updateProject(w http.ResponseWriter, r *http.Request) {

}

func deleteProject(w http.ResponseWriter, r *http.Request) {

}

func getColumns(w http.ResponseWriter, r *http.Request) {

}

func createColumn(w http.ResponseWriter, r *http.Request) {

}

func updateColumn(w http.ResponseWriter, r *http.Request) {

}

func deleteColumn(w http.ResponseWriter, r *http.Request) {

}

func swapColumns(w http.ResponseWriter, r *http.Request) {

}
