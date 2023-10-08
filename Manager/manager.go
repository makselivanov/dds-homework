package manager

import (
	database "Database"
	"log"
	"time"
)

var clock = make(map[string]uint64)
var localSource string

var db database.Database

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

func AddTransaction(transaction database.Transaction) {
	db.AddTransaction(transaction)
}

func GetValue() string {
	return db.GetValue()
}

func Init(source string) {
	localSource = source
	clock[localSource] = 0
	db = database.NewDatabase(clock)
	go autoSaveSnapshot()
	database.Init()
}

func GetVClock() map[string]uint64 {
	return clock
}
