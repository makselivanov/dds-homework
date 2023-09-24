package database

import "log"

var globalVersion int = 0

type Snapshot struct {
	value   string
	version int
}

type Transaction struct {
	fromSnapshotVersion int
	value               string
}

type Database struct {
	snapshot     Snapshot
	transactions []Transaction
}

func NewSnapshot(newValue string) Snapshot {
	globalVersion++
	return Snapshot{
		value:   newValue,
		version: globalVersion,
	}
}

func NewDatabase() Database {
	log.Println("Create database")
	return Database{
		snapshot:     NewSnapshot(""),
		transactions: make([]Transaction, 0),
	}
}

func (db *Database) AddTransaction(value string) {
	curVersion := db.snapshot.version
	transaction := Transaction{
		fromSnapshotVersion: curVersion,
		value:               value,
	}
	//FIXME not thread safe?
	db.transactions = append(db.transactions, transaction)
	log.Printf("Add new transaction from snapshot version %d\n", curVersion)
}

func (db Database) GetValue() string {
	//FIXME should be tread safe?
	shapshot := db.snapshot
	value := shapshot.value
	if len(db.transactions) > 0 && db.transactions[len(db.transactions)-1].fromSnapshotVersion >= shapshot.version {
		index := len(db.transactions) - 1
		transaction := db.transactions[index]
		value = transaction.value
		log.Printf("Last value from transaction version from snapshot v%d and index %d", transaction.fromSnapshotVersion, index)
	} else {
		log.Printf("Last value from snapshot version %d", shapshot.version)
	}
	log.Printf("Return value from database: %s\n", value)
	return value
}

func (db *Database) SaveSnapshot() {
	if len(db.transactions) == 0 {
		return
	}
	operation := db.transactions[len(db.transactions)-1]
	if operation.fromSnapshotVersion >= globalVersion {
		db.snapshot = NewSnapshot(operation.value)
		log.Printf("Save new snapshot version %d\n", db.snapshot.version)
	}
}
