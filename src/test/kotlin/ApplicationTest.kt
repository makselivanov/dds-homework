import io.ktor.client.request.*
import io.ktor.client.statement.*
import io.ktor.http.*
import io.ktor.server.testing.*
import io.ktor.util.*
import kotlin.test.Test
import kotlin.test.assertContentEquals
import kotlin.test.assertEquals

class ApplicationTest {
    @Test
    fun testPutGetDelete() = testApplication {
        val message = "Hello, internet world!!!"
        var response = client.put("/replace") {
            setBody(message)
        }
        assertEquals(HttpStatusCode.OK, response.status)
        response = client.get("/get")
        assertEquals(HttpStatusCode.OK, response.status)
        assertContentEquals(message.toByteArray(), response.bodyAsChannel().toByteArray())
        response = client.delete("/replace")
        assertEquals(HttpStatusCode.OK, response.status)
        response = client.get("/get")
        assertEquals(HttpStatusCode.OK, response.status)
        assertContentEquals(ByteArray(0), response.bodyAsChannel().toByteArray())

    }
}