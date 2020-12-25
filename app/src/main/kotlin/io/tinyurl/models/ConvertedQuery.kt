package io.tinyurl.models

data class ConvertedQuery(
    val status: Boolean,
    val q: String,
    val result: String,
    val stored: Boolean
)
