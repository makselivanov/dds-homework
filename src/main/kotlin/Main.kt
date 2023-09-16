import io.ktor.server.application.*
import io.ktor.server.netty.*
import plugins.configureRouting

fun Application.module() {
    configureRouting()
}


fun main(args: Array<String>) {
    EngineMain.main(args)
}