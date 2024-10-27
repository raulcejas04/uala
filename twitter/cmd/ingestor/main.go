package main

import (
	"fmt"
	"net/http"

	kafka "twitter/pkg/kafka"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	TOPIC = "services"
)

func channelHandlerPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		usuario := vars["usuario"]
		mensaje := vars["mensaje"]
		log.Infof("Nuevo twitter %s %s", usuario, mensaje)
		kafka.Produce(TOPIC, fmt.Sprintf("%s;%s", usuario, mensaje))

	}
}

func channelHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		usuario := vars["usuario"]
		log.Infof("Nuevo twitter %s %s", usuario)
	}
}

func handleRequests() {
	r := mux.NewRouter()

	chPost := channelHandlerPost()
	chGet := channelHandlerGet()

	r.HandleFunc("/{usuario}/{mensaje}", chPost).Methods("POST")
	r.HandleFunc("/{id}", chGet).Methods("GET")

	log.Infof("Inicializando servidor ")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	srv.ListenAndServe()

}

func main() {

	handleRequests()

}
