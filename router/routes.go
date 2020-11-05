package routes

import (
	"github.com/gorilla/mux"
)

func Routes(r *mux.Router) {
	u := r.PathPrefix("/user").Subrouter()
	userRoutes(u)
	p := r.PathPrefix("/patient").Subrouter()
	patientRoutes(p)
	res := r.PathPrefix("/exams").Subrouter()
	examsRoutes(res)
	f := r.PathPrefix("/file").Subrouter()
	fileRoutes(f)
}