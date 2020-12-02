package routes

import (
	"github.com/gorilla/mux"
)

func Routes(r *mux.Router) {
	u := r.PathPrefix("/systemuser").Subrouter()
	userRoutes(u)
	p := r.PathPrefix("/patient").Subrouter()
	patientRoutes(p)
	res := r.PathPrefix("/exams").Subrouter()
	examsRoutes(res)
	f := r.PathPrefix("/file").Subrouter()
	fileRoutes(f)
	prt := r.PathPrefix("/protocol").Subrouter()
	protocolRoutes(prt)
	l := r.PathPrefix("/location").Subrouter()
	locationRoutes(l)
}