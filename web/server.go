package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func StartServer() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	var db db.Database = &db.SqlDb{}
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

}

func getProject(w http.ResponseWriter, r *http.Request) {

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
