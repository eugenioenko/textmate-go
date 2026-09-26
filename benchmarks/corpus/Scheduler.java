package com.example.scheduling;

import java.time.Duration;
import java.time.Instant;
import java.util.ArrayList;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Supplier;

/**
 * Schedules named tasks with retry and backoff.
 *
 * @param <T> the task result type
 */
public final class Scheduler<T> implements AutoCloseable {
    private static final int MAX_ATTEMPTS = 5;
    private static final Duration BASE_DELAY = Duration.ofMillis(250);

    private final ConcurrentHashMap<String, Task<T>> tasks = new ConcurrentHashMap<>();
    private volatile boolean closed;

    public record Task<T>(String name, Supplier<T> body, Instant createdAt) {
        public Task {
            Objects.requireNonNull(name, "name");
            if (name.isBlank()) {
                throw new IllegalArgumentException("task name must not be blank");
            }
        }
    }

    public sealed interface Outcome<T> permits Success, Failure {}

    public record Success<T>(T value, int attempts) implements Outcome<T> {}

    public record Failure<T>(Exception cause, int attempts) implements Outcome<T> {}

    public void submit(String name, Supplier<T> body) {
        if (closed) {
            throw new IllegalStateException("scheduler closed");
        }
        tasks.putIfAbsent(name, new Task<>(name, body, Instant.now()));
    }

    public Optional<Outcome<T>> run(String name) throws InterruptedException {
        Task<T> task = tasks.get(name);
        if (task == null) {
            return Optional.empty();
        }
        Exception last = null;
        for (int attempt = 1; attempt <= MAX_ATTEMPTS; attempt++) {
            try {
                return Optional.of(new Success<>(task.body().get(), attempt));
            } catch (RuntimeException e) {
                last = e;
                Thread.sleep(BASE_DELAY.multipliedBy(1L << (attempt - 1)).toMillis());
            }
        }
        return Optional.of(new Failure<>(last, MAX_ATTEMPTS));
    }

    public List<String> pending() {
        var names = new ArrayList<>(tasks.keySet());
        names.sort(String::compareToIgnoreCase);
        return List.copyOf(names);
    }

    @Override
    public void close() {
        closed = true;
        tasks.clear();
    }
}
