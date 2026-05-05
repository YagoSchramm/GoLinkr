package router

import "github.com/gorilla/mux"

func Mount(root *mux.Router, modules ...Module) {
	for _, module := range modules {
		subrouter := root.PathPrefix(module.Path()).Subrouter()
		for _, middleware := range module.Middlewares() {
			subrouter.Use(middleware)
		}
		for _, route := range module.Routes() {
			subrouter.HandleFunc(route.Path, route.Handler).Methods(route.HttpMethods...)
		}
	}
}
