package main

import (
	manager "Manager"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var source = "selivanov"

var peers []string

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
		transaction := manager.NewTransaction(string(buffer[:n]), source)
		manager.AddTransaction(transaction)
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
	log.Println("GET /test")
	http.ServeFile(writer, reader, "static/index.html")
}

func requestVClock(writer http.ResponseWriter, reader *http.Request) {
	switch reader.Method {
	case http.MethodGet:
		log.Println("GET /ws")
		writer.Header().Set("Content-Type", "application/json")

		jsonStr, err := json.Marshal(manager.GetVClock())
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
		} else {
			writer.WriteHeader(http.StatusOK)
			writer.Write(jsonStr)
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("Should add program args with bind port and hostpeers")
	}
	port := os.Args[1]
	peers = os.Args[2:]
	manager.Init(source)
	log.Printf("Binding to port %s", port)
	http.HandleFunc("/replace", requestReplace) // Устанавливаем роутер
	http.HandleFunc("/post", requestReplace)
	http.HandleFunc("/get", requestGet)
	http.HandleFunc("/test", requestTest)
	http.HandleFunc("/vclock", requestVClock)
	//http.HandleFunc("/ws", )
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil) // устанавливаем порт веб-сервера
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
