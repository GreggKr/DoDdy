package com.github.dod.doddy.core

import net.dv8tion.jda.core.entities.Member
import net.dv8tion.jda.core.events.message.MessageReceivedEvent
import java.lang.Exception
import kotlin.reflect.KFunction
import kotlin.reflect.KParameter
import kotlin.reflect.jvm.javaType

data class CommandFunction(
    val module: Module,
    val function: KFunction<*>,
    val parameters: List<KParameter>,
    val optionals: List<Int>,
    val allArgs: Boolean,
    val commandAnnotation: Command
) {

    companion object {
        private val snowflakeRegex = Regex("\\d{17,19}")
        private val mentionRegex = Regex("<@!?\\d{17,19}>")
        private val usernameDiscrimRegex = Regex(".+#\\d{4}")

        private val stringType = String::class.java
        private val intType = Int::class.java
        private val longType = Long::class.java
        private val shortType = Short::class.java
        private val doubleType = Double::class.java
        private val memberType = Member::class.java
    }

    suspend fun call(event: MessageReceivedEvent, args: List<String>): CommandResult {
        if (args.size + optionals.size < parameters.size && !allArgs) {//TODO: check for too many arguments
            return InvalidArgs(args)
        }
        val params = ArrayList<Any?>()
        params.add(module)
        params.add(event)
        if (allArgs) {
            params.addAll(args)
        }
        args.forEachIndexed { index, argument ->
            val paramIndex = index + 2
            if (parameters.size > index) {
                when (parameters[index].type.javaType) {
                    stringType -> {
                        params.add(paramIndex, argument)
                    }
                    intType -> {
                        val number = argument.toIntOrNull()
                        if (number != null) {
                            params.add(paramIndex, number)
                        } else {
                            return InvalidArg(argument, "not a number")
                        }
                    }
                    longType -> {
                        val number = argument.toLongOrNull()
                        if (number != null) {
                            params.add(paramIndex, number)
                        } else {
                            return InvalidArg(argument, "not a number")
                        }
                    }
                    shortType -> {
                        val number = argument.toShortOrNull()
                        if (number != null) {
                            params.add(paramIndex, number)
                        } else {
                            return InvalidArg(argument, "not a number")
                        }
                    }
                    doubleType -> {
                        val number = argument.toDoubleOrNull()
                        if (number != null) {
                            params.add(paramIndex, number)
                        } else {
                            return InvalidArg(argument, "not a number")
                        }
                    }
                    memberType -> {
                        val member: Member? = when {
                            snowflakeRegex.matches(argument) -> event.guild.getMemberById(argument)
                            mentionRegex.matches(argument) -> event.guild.getMemberById(argument.slice(2 until argument.length - 1).removePrefix("!"))
                            usernameDiscrimRegex.matches(argument) -> {
                                val hashIndex = argument.lastIndexOf("#")
                                val username = argument.slice(0 until hashIndex)
                                val discrim = argument.substring(hashIndex + 1)
                                val mem = event.guild.getMembersByName(username, true).filter {
                                    it.user.discriminator == discrim
                                }
                                mem.getOrNull(0)
                            }
                            else -> {
                                val mem1 = event.guild.getMembersByNickname(argument, true)
                                val mem2 = event.guild.getMembersByName(argument, true)
                                mem1.getOrNull(0) ?: mem2.getOrNull(0)
                            }
                        }
                        if (member != null) {
                            params.add(paramIndex, member)
                        } else {
                            return InvalidArg(argument, "not a known member in this guild")
                        }
                    }
                }
            }
        }
        optionals.forEach {
            if (params.size <= it || params[it] == null) {
                params.add(it, null)
            }
        }
        try {
            function.call(*params.toArray())
        } catch (exception: Exception) {
            return CommandError(exception)
        }
        return Success("bla")
    }
}