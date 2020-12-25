package io.tinyurl.repositories

import io.ktor.client.HttpClient
import io.ktor.client.request.get
import io.tinyurl.models.ConvertedQuery
import kotlinx.coroutines.async
import kotlinx.coroutines.coroutineScope

interface TinyUrlRepository {
    suspend fun convert(originUrl: String): ConvertedQuery
}

class TinyUrlRepositoryImpl(private val client: HttpClient) : TinyUrlRepository {

    override suspend fun convert(originUrl: String) = coroutineScope {
        val request =
            async { client.get<ByteArray>("http://tinyurl.com/api-create.php?url=$originUrl") }
        val resultBytes = request.await()

        client.close()
        return@coroutineScope ConvertedQuery(
            true,
            originUrl,
            String(resultBytes),
            false,
        )
    }
}