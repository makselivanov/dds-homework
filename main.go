package main

import (
	database "Database"
	manager "Manager"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var source string

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
		writer.Header().Set("Content-Type", "application/json")
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
		log.Println("GET /vclock")
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

func websocketHandler(writer http.ResponseWriter, reader *http.Request) {
	ws, err := websocket.Accept(writer, reader, &websocket.AcceptOptions{InsecureSkipVerify: true, OriginPatterns: []string{"*"}})
	if err != nil {
		log.Printf("Problem during accept\n")
		return
	}
	defer ws.Close(websocket.StatusInternalError, fmt.Sprintf("Connection is closed with %s", source))
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ch := make(chan database.Transaction)

	manager.SendAllTransactions(ch)

loop:
	for {
		select {
		case <-ctx.Done():
			return
		case transaction, ok := <-ch:
			if !ok {
				break loop
			}
			log.Printf("Sending transaction id %d source %s to peer\n", transaction.Id, transaction.Source)
			wsjson.Write(ctx, ws, transaction)
		}
	}

	ch = make(chan database.Transaction)

	manager.AddChannel(ch)

	for {
		select {
		case <-ctx.Done():
			return
		case transaction := <-ch:
			log.Printf("Sending transaction id %d source %s to peer\n", transaction.Id, transaction.Source)
			wsjson.Write(ctx, ws, transaction)
		}
	}
}

func runLoop(peer string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://%s/ws", peer), nil)
	if err != nil {
		log.Printf("Problem connecting with dial to %s", peer)
		return
	}
	log.Printf("Connected with %s", peer)
	defer c.Close(websocket.StatusInternalError, fmt.Sprintf("Connection is closed with %s", source))
	var transaction database.Transaction

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			err := wsjson.Read(ctx, c, &transaction)
			if err != nil {
				//log.Printf("Error read from %s: %s", peer, err.Error())
				continue loop
			}
			log.Printf("Got transaction from %s with id %d source %s", peer, transaction.Id, transaction.Source)
			manager.AddTransaction(transaction)
		}
	}
}

func runPeer(peer string) {
	time.Sleep(time.Second * 2)
	for {
		runLoop(peer)
		log.Printf("Connection lost with %s, Reconnecting...", peer)
		time.Sleep(time.Second * 5)
	}
}

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("Should add program args with bind port and hostpeers")
	}
	source = os.Args[1]
	port := os.Args[2]
	peers = os.Args[3:]
	manager.Init(source)
	log.Printf("Binding to port %s", port)
	http.HandleFunc("/replace", requestReplace) // Устанавливаем роутер
	http.HandleFunc("/post", requestReplace)
	http.HandleFunc("/get", requestGet)
	http.HandleFunc("/test", requestTest)
	http.HandleFunc("/vclock", requestVClock)
	http.HandleFunc("/ws", websocketHandler)

	for _, peer := range peers {
		go runPeer(peer)
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil) // устанавливаем порт веб-сервера
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
