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
	"github.com/i-hate-nicknames/tasker/tasker"
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

type projectRequest struct {
	Name        string
	Description string
}

func createProject(w http.ResponseWriter, r *http.Request) {
	var req projectRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// todo: add validation library
	if req.Name == "" || req.Description == "" {
		http.Error(w, "Invalid request structure", http.StatusBadRequest)
		return
	}
	db := r.Context().Value(keyDb).(db.Database)
	// todo: probably move all this stuff to service layer, which should call db on its own
	project := tasker.MakeProject(nil, req.Name, req.Description)
	err = db.SaveProject(project)
	if err != nil {
		http.Error(w, "Failed to save project", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("ok"))
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.Context().Value(keyProjectId).(int)
	db := r.Context().Value(keyDb).(db.Database)
	project, err := db.GetProject(projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("project %d not found", projectID), http.StatusNotFound)
		return
	}

	var req projectRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// todo: add validation library
	if req.Name == "" || req.Description == "" {
		http.Error(w, "Invalid request structure", http.StatusBadRequest)
		return
	}
	project.Name = req.Name
	project.Description = req.Description
	err = db.SaveProject(project)
	if err != nil {
		http.Error(w, "Failed to save project", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("ok"))
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
		log.Printf("Error encoding to json %v, err: %s\n", marshaled, err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	w.Write(marshaled)
}
