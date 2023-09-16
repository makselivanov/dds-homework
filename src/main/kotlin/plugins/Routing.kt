package plugins

import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import storage.Storage

val storage = Storage()
fun Application.configureRouting() {
    routing {
        route("/replace") {
            post {
                val body = call.receive<ByteArray>()
                storage.update(body)
                call.respond(HttpStatusCode.OK)
            }
            put {
                val body = call.receive<ByteArray>()
                storage.update(body)
                call.respond(HttpStatusCode.OK)
            }
            delete {
                storage.update(ByteArray(0))
                call.respond(HttpStatusCode.OK)
            }
        }
        route("/get") {
            get {
                val body = storage.get()
                call.respond(HttpStatusCode.OK, body)
            }
        }
    }
}