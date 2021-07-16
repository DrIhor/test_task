package routes

import (
	"github.com/gorilla/mux"
)

func Handler() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/iteams", showAllIteams).Methods("GET")
	router.HandleFunc("/iteams/{name}", showAllIteams).Methods("GET")
	router.HandleFunc("/iteams", addNewIteam).Methods("POST")
	router.HandleFunc("/iteams/{name}", buyIteams).Methods("PUT")
	router.HandleFunc("/iteams/{name}", deleteIteam).Methods("DELETE")

	return router
}
