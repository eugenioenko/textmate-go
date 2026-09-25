final class Example {
    /* A comment whose end is
       on another line. */
    static String message(String name) {
        return """
            hello,
            %s
            """.formatted(name);
    }
}
