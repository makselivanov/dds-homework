package Database

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

func newSnapshot(newValue string) Snapshot {
	globalVersion++
	return Snapshot{
		value:   newValue,
		version: globalVersion - 1,
	}
}

func newDatabase() Database {
	return Database{
		snapshot:     newSnapshot(""),
		transactions: make([]Transaction, 0),
	}
}
