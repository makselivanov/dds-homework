package storage

class Storage {
    private var array = ByteArray(0)

    fun update(newData: ByteArray) {
        array = newData.copyOf()
    }
    fun get(): ByteArray {
        return array
    }
}