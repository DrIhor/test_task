package routes

import (
	"net/http"

	"github.com/DrIhor/test_task/internal/models/iteams"
	servIteams "github.com/DrIhor/test_task/internal/service/iteams"

	"github.com/gorilla/mux"
)

func Handler(router *mux.Router, stor iteams.IteamServices) {
	router.HandleFunc("/iteams", func(w http.ResponseWriter, r *http.Request) {
		servIteams.ShowAllIteams(w, r, stor)
	}).Methods("GET")
	router.HandleFunc("/iteams/{name}", func(w http.ResponseWriter, r *http.Request) {
		servIteams.ShowIteam(w, r, stor)
	}).Methods("GET")
	router.HandleFunc("/iteams", func(w http.ResponseWriter, r *http.Request) {
		servIteams.AddNewIteam(w, r, stor)
	}).Methods("POST")
	router.HandleFunc("/iteams/{name}", func(w http.ResponseWriter, r *http.Request) {
		servIteams.BuyIteams(w, r, stor)
	}).Methods("PUT")
	router.HandleFunc("/iteams/{name}", func(w http.ResponseWriter, r *http.Request) {
		servIteams.DeleteIteam(w, r, stor)
	}).Methods("DELETE")
}
