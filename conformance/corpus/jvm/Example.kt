data class Person(val name: String)

/* A Kotlin comment
   crosses a line. */
fun greet(person: Person): String {
    val detail = "${person.name.length} letters"
    return "Hello ${person.name} ($detail)"
}
