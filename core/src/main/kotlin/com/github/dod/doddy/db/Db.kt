package com.github.dod.doddy.db

import org.litote.kmongo.async.KMongo

object Db {
    val client = KMongo.createClient()
    val instance = client.getDatabase("mongo")
}