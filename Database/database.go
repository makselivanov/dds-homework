package database

import (
	"log"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

type Snapshot struct {
	snap string
}

type Transaction struct {
	Source  string
	Id      uint64
	Payload string
}

type Database struct {
	snapshot     Snapshot
	transactions []Transaction
}

var optionsJsonPatch *jsonpatch.ApplyOptions

func Init() {
	optionsJsonPatch = jsonpatch.NewApplyOptions()
	optionsJsonPatch.AllowMissingPathOnRemove = true
}

func (db Database) GetTransactions() []Transaction {
	return db.transactions
}

func NewSnapshot(newValue string) Snapshot {
	return Snapshot{
		snap: newValue,
	}
}

func NewDatabase(clock map[string]uint64) Database {
	log.Println("Create database")
	return Database{
		snapshot:     NewSnapshot("{}"),
		transactions: make([]Transaction, 0),
	}
}

func (db *Database) AddTransaction(transaction Transaction) {
	log.Printf("add transaction id %d source %s", transaction.Id, transaction.Source)
	db.transactions = append(db.transactions, transaction)
}

func ApplyTransaction(snap string, transaction Transaction) (string, error) {
	patch, err := jsonpatch.DecodePatch([]byte(transaction.Payload))
	if err != nil {
		log.Printf("Error when trying to decode patch from %s id %d\n", transaction.Source, transaction.Id)
		return snap, err
	}
	newsnap, err := patch.ApplyWithOptions([]byte(snap), optionsJsonPatch)
	if err != nil {
		log.Printf("Error when trying to apply patch from %s id %d\n", transaction.Source, transaction.Id)
		return snap, err
	}
	return string(newsnap), nil
}

func (db Database) GetValue() string {
	//FIXME should be tread safe?
	db.SaveSnapshot()
	snapshot := db.snapshot
	snap := snapshot.snap
	return snap
}

func (db *Database) SaveSnapshot() {
	snapshot := db.snapshot
	newsnap := snapshot.snap
	transactions := db.transactions
	flag := false
	var err error = nil
	for _, transaction := range transactions {
		newsnap, err = ApplyTransaction(newsnap, transaction)
		if err == nil {
			flag = true
		}
	}
	if flag {
		log.Printf("Collect new snapshot\n")
		db.snapshot = NewSnapshot(newsnap)
	} else {
		log.Printf("Nothing change outside of snapshot")
	}
}
