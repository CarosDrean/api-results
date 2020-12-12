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
	prtsu := r.PathPrefix("/protocolsystemuser").Subrouter()
	protocolSystemUserRoutes(prtsu)
	l := r.PathPrefix("/location").Subrouter()
	locationRoutes(l)
	pro := r.PathPrefix("/protocol").Subrouter()
	protocolRoutes(pro)
	s := r.PathPrefix("/service").Subrouter()
	serviceRoutes(s)
	o := r.PathPrefix("/organization").Subrouter()
	organizationRoutes(o)
	sp := r.PathPrefix("/system-parameter").Subrouter()
	systemParameterRoutes(sp)
	cmp := r.PathPrefix("/component").Subrouter()
	componentRoutes(cmp)
	st := r.PathPrefix("/statistic").Subrouter()
	statisticRoutes(st)
}