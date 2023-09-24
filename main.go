package main

import (
	database "Database"
	"io"
	"log"
	"net/http"
	"time"
)

var db = database.NewDatabase()

func requestReplace(writer http.ResponseWriter, reader *http.Request) {
	switch reader.Method {
	case http.MethodPost:
		log.Println("POST /replace")
		var buffer string
		_, err := reader.Body.Read([]byte(buffer))

		if err != nil && err != io.EOF {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
		db.AddTransaction(buffer)
	default:
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func requestGet(writer http.ResponseWriter, reader *http.Request) {
	switch reader.Method {
	case http.MethodGet:
		log.Println("GET /get")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(db.GetValue()))
	default:
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func autoSaveSnapshot() {
	for {
		time.Sleep(time.Minute)
		db.SaveSnapshot()
	}
}

func main() {
	go autoSaveSnapshot()
	http.HandleFunc("/replace", requestReplace) // Устанавливаем роутер
	http.HandleFunc("/get", requestGet)
	err := http.ListenAndServe(":8080", nil) // устанавливаем порт веб-сервера
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
