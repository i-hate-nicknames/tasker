package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/i-hate-nicknames/tasker/db"
)

type ctxKey int

const (
	keyDb ctxKey = iota
	keyProjectId
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
			r.Use(intURIParamMiddleware("projectID", keyProjectId))
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
	http.ListenAndServe(":3000", r)
}

func intURIParamMiddleware(paramName string, key interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			val := chi.URLParam(r, paramName)
			intval, err := strconv.Atoi(val)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s must be an int", paramName), http.StatusBadRequest)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), key, intval))
			next.ServeHTTP(w, r)
		})
	}
}

func listProjects(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(keyDb).(db.Database)
	projects, err := db.GetProjects()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	writeJson(w, projects)
}

func getProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.Context().Value(keyProjectId).(int)
	db := r.Context().Value(keyDb).(db.Database)
	project, err := db.GetProject(projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("project %d not found", projectID), http.StatusNotFound)
		return
	}
	writeJson(w, project)
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

func writeJson(w http.ResponseWriter, val interface{}) {
	marshaled, err := json.Marshal(val)
	if err != nil {
		log.Printf("Error encoding to json %v\n", marshaled)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	w.Write(marshaled)
}
