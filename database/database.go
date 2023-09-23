package database

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
		version: globalVersion - 1,
	}
}

func NewDatabase() Database {
	return Database{
		snapshot:     NewSnapshot(""),
		transactions: make([]Transaction, 0),
	}
}
