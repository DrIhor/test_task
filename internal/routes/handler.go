package routes

import "github.com/gorilla/mux"

func postgresRoutes(r *mux.Router) *mux.Router {
	return r
}

func mongoRoutes(r *mux.Router) *mux.Router {
	return r
}

func Handler() *mux.Router {
	router := mux.NewRouter()

	// postgres
	postgresRoutes(router.PathPrefix("/postgres").Subrouter())
	// mongo
	mongoRoutes(router.PathPrefix("/mongo").Subrouter())
	return router
}
