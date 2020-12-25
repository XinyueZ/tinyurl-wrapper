/*
 * Copyright 2020 by Chris Xinyue
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package io.tinyurl

import io.ktor.application.Application
import io.ktor.application.call
import io.ktor.application.install
import io.ktor.client.HttpClient
import io.ktor.client.engine.cio.CIO
import io.ktor.features.CallLogging
import io.ktor.features.ContentNegotiation
import io.ktor.features.DefaultHeaders
import io.ktor.gson.gson
import io.ktor.http.HttpStatusCode
import io.ktor.response.respond
import io.ktor.routing.get
import io.ktor.routing.routing
import io.tinyurl.repositories.TinyUrlRepositoryImpl
import io.tinyurl.viewmodels.TinyUrlViewModel

private val String?.isUri: Boolean
    get() {
        return if (this.isNullOrBlank()) false
        else {
            val pattern =
                """^((((H|h)(T|t)|(F|f))(T|t)(P|p)((S|s)?))\://)?(www.|[a-zA-Z0-9].)[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,6}(\:[0-9]{1,5})*(/($|[a-zA-Z0-9\.\,\;\?\'\\\+&amp;%\$#\=~_\-]+))*$"""
            pattern.toRegex().matches(this)
        }
    }

fun Application.main() {
    install(DefaultHeaders)
    install(CallLogging)
    install(ContentNegotiation) {
        gson {
            this.setPrettyPrinting()
        }
    }

    val vm = TinyUrlViewModel(TinyUrlRepositoryImpl(HttpClient(CIO)))

    routing {
        get("/") {
            if (call.request.queryParameters.isEmpty() || !call.request.queryParameters["q"].isUri) {
                call.respond(HttpStatusCode.NoContent)
            } else {
                call.respond(
                    HttpStatusCode.OK,
                    vm.convert(call.request.queryParameters["q"]!!),
                )
            }
        }
    }
}