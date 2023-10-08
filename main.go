package main

import (
	databaseManager "Manager"
	manager "Manager"
	"io"
	"log"
	"net/http"
)

func requestReplace(writer http.ResponseWriter, reader *http.Request) {
	switch reader.Method {
	case http.MethodPost:
		var buffer []byte = make([]byte, 1000)
		n, err := reader.Body.Read(buffer)
		log.Printf("POST /replace body len: %d", n)

		if err != nil && err != io.EOF {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
		manager.AddTransaction(string(buffer[:n]))
	default:
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func requestGet(writer http.ResponseWriter, reader *http.Request) {
	switch reader.Method {
	case http.MethodGet:
		log.Println("GET /get")
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "text/plain")
		writer.Write([]byte(manager.GetValue()))
	default:
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func requestTest(writer http.ResponseWriter, reader *http.Request) {
	http.ServeFile(writer, reader, "static/index.html")
}

func main() {
	databaseManager.Init()
	http.HandleFunc("/replace", requestReplace) // Устанавливаем роутер
	http.HandleFunc("/post", requestReplace)
	http.HandleFunc("/get", requestGet)
	http.HandleFunc("/test", requestTest)
	//http.HandleFunc("/vclock", )
	//http.HandleFunc("/ws", )
	err := http.ListenAndServe(":8080", nil) // устанавливаем порт веб-сервера
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
