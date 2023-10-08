package manager

import (
	database "Database"
	"log"
	"time"
)

var db = database.NewDatabase()

func autoSaveSnapshot() {
	for {
		time.Sleep(time.Minute)
		log.Println("Trying to save snapshot")
		db.SaveSnapshot()
	}
}

func AddTransaction(newValue string) {
	db.AddTransaction(newValue)
}

func GetValue() string {
	return db.GetValue()
}

func Init() {
	go autoSaveSnapshot()
}
