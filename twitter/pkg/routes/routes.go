package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"twitter/pkg/kafka"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Response struct {
	Message string `json:"message"`
}

func ChannelHandlerPost(topic string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fmt.Printf("parameters %+v'n", vars)
		usuario := vars["usuario"]
		mensaje := vars["mensaje"]
		respuesta := fmt.Sprintf("Nuevo twitter %s %s", usuario, mensaje)
		kafka.Produce(topic, fmt.Sprintf("%s;%s", usuario, mensaje))
		response := Response{Message: respuesta}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	}
}

func ChannelHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		usuario := vars["usuario"]
		log.Infof("Nuevo twitter %s ", usuario)
	}
}

func HandleRequests(topic string) {
	r := mux.NewRouter()

	chPost := ChannelHandlerPost(topic)
	chGet := ChannelHandlerGet()

	r.HandleFunc("/{usuario}/{mensaje}", chPost).Methods("POST")
	r.HandleFunc("/{id}", chGet).Methods("GET")

	log.Infof("Inicializando servidor ")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	srv.ListenAndServe()

}
