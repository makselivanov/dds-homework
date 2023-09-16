package storage

class Storage {
    private val array = ByteArray(0)

    fun update(newData: ByteArray) {
        newData.copyInto(array)
    }
    fun get(): ByteArray {
        return array
    }
}