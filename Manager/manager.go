package manager

import (
	database "Database"
	"log"
	"time"
)

var clock = make(map[string]uint64)
var localSource string

var db database.Database
var channels []chan database.Transaction
var channel chan database.Transaction

func NewTransaction(payload string, source string) database.Transaction {
	clock[localSource]++
	return database.Transaction{
		Id:      clock[localSource],
		Payload: payload,
		Source:  source,
	}
}

func autoSaveSnapshot() {
	for {
		time.Sleep(time.Second * 30)
		log.Println("Trying to save snapshot")
		db.SaveSnapshot()
	}
}

func autoSendTransactions() {
	for {
		var transaction = <-channel
		//check if valid?

		if transaction.Id > clock[transaction.Source] {
			clock[transaction.Source] = transaction.Id

			for _, ch := range channels {
				ch <- transaction
			}
		}

	}
}

func AddChannel(ch chan database.Transaction) {
	channels = append(channels, ch)
}

func autoLocalSave() {
	ch := make(chan database.Transaction)
	channels = append(channels, ch)

	for {
		var transaction = <-ch
		db.AddTransaction(transaction)
	}
}

func AddTransaction(transaction database.Transaction) {
	channel <- transaction
}

func GetValue() string {
	return db.GetValue()
}

func Init(source string) {
	localSource = source
	channels = make([]chan database.Transaction, 0)
	channel = make(chan database.Transaction)
	clock[localSource] = 0
	db = database.NewDatabase(clock)
	go autoSaveSnapshot()
	go autoSendTransactions()
	go autoLocalSave()
	database.Init()
}

func GetVClock() map[string]uint64 {
	return clock
}
