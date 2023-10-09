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
var localVersion uint64 = 0

func NewTransaction(payload string, source string) database.Transaction {
	localVersion++
	log.Printf("creating new transaction id %d from %s\n", localVersion, source)
	return database.Transaction{
		Id:      localVersion,
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
		if transaction.Source == localSource {
			localVersion = max(transaction.Id, localVersion)
		}
		if transaction.Id > clock[transaction.Source] {
			clock[transaction.Source] = transaction.Id
			log.Printf("Got new transaction id %d from %s", transaction.Id, transaction.Source)
			for _, ch := range channels {
				ch <- transaction
			}
		} else {
			log.Printf("Got old transaction id %d from %s, local version %d", transaction.Id, transaction.Source, clock[transaction.Source])
		}

	}
}

func AddChannel(ch chan database.Transaction) {
	channels = append(channels, ch)
}

func sendingAllTransactions(ch chan database.Transaction) {
	for _, transaction := range db.GetTransactions() {
		ch <- transaction
	}
	close(ch)
}

func SendAllTransactions(ch chan database.Transaction) {
	go sendingAllTransactions(ch)
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
