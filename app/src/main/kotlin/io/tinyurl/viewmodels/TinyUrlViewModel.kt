package io.tinyurl.viewmodels

import io.tinyurl.models.ConvertedQuery
import io.tinyurl.repositories.TinyUrlRepository

class TinyUrlViewModel(private val tinyUrlRepository: TinyUrlRepository) {
    suspend fun convert(originUrl: String): ConvertedQuery {
        return tinyUrlRepository.convert(originUrl)
    }
}