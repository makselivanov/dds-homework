package database

import (
	"log"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

type Snapshot struct {
	snap  string
	clock map[string]uint64
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

func NewSnapshot(newValue string, clock map[string]uint64) Snapshot {
	newclock := make(map[string]uint64)
	for key, value := range clock {
		newclock[key] = value
	}
	return Snapshot{
		snap:  newValue,
		clock: newclock,
	}
}

func NewDatabase(clock map[string]uint64) Database {
	log.Println("Create database")
	return Database{
		snapshot:     NewSnapshot("{}", clock),
		transactions: make([]Transaction, 0),
	}
}

func (db *Database) AddTransaction(transaction Transaction) {
	db.transactions = append(db.transactions, transaction)
}

func ApplyTransaction(snap string, transaction Transaction) (string, error) {
	patch, err := jsonpatch.DecodePatch([]byte(transaction.Payload))
	if err != nil {
		log.Println("Error when trying to decode patch")
		return snap, err
	}
	newsnap, err := patch.ApplyWithOptions([]byte(snap), optionsJsonPatch)
	if err != nil {
		log.Println("Error when trying to apply patch")
		return snap, err
	}
	return string(newsnap), nil
}

func (db Database) GetValue() string {
	//FIXME should be tread safe?
	snapshot := db.snapshot
	snap := snapshot.snap
	return snap
}

func (db *Database) SaveSnapshot() {
	snapshot := db.snapshot
	newsnap := snapshot.snap
	transactions := db.transactions
	flag := false
	newclock := snapshot.clock
	var err error = nil
	for _, transcation := range transactions {
		if transcation.Id > snapshot.clock[transcation.Source] {
			newsnap, err = ApplyTransaction(newsnap, transcation)
			if err == nil {
				flag = true
				newclock[transcation.Source] = max(newclock[transcation.Source], transcation.Id)
			}
		}
	}
	if flag {
		db.snapshot = NewSnapshot(newsnap, newclock)
	}
}
