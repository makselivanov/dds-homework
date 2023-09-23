package main

import (
	database "Database"
	"log"
	"net/http"
)

var db = database.NewDatabase()

func requestReplace(writer http.ResponseWriter, reader *http.Request) {
	switch reader.Method {
	case http.MethodPost:
		//TODO
	default:
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func requestGet(writer http.ResponseWriter, reader *http.Request) {
	switch reader.Method {
	case http.MethodGet:
		//TODO
	default:
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/replace", requestReplace) // Устанавливаем роутер
	http.HandleFunc("/get", requestGet)
	err := http.ListenAndServe(":8080", nil) // устанавливаем порт веб-сервера
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
