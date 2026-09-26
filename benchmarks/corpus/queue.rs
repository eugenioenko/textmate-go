//! A bounded, priority-aware job queue.

use std::cmp::Ordering;
use std::collections::BinaryHeap;
use std::fmt;
use std::time::{Duration, Instant};

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum Priority {
    Low,
    Normal,
    High(u8),
}

#[derive(Debug)]
pub struct Job<'a> {
    pub id: u64,
    pub name: &'a str,
    pub priority: Priority,
    pub enqueued: Instant,
}

impl<'a> Ord for Job<'a> {
    fn cmp(&self, other: &Self) -> Ordering {
        let rank = |p: &Priority| match p {
            Priority::Low => 0u16,
            Priority::Normal => 1,
            Priority::High(level) => 2 + *level as u16,
        };
        rank(&self.priority)
            .cmp(&rank(&other.priority))
            .then_with(|| other.enqueued.cmp(&self.enqueued))
    }
}

impl<'a> PartialOrd for Job<'a> {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

impl<'a> PartialEq for Job<'a> {
    fn eq(&self, other: &Self) -> bool {
        self.id == other.id
    }
}

impl<'a> Eq for Job<'a> {}

pub struct Queue<'a> {
    heap: BinaryHeap<Job<'a>>,
    capacity: usize,
}

#[derive(Debug)]
pub struct QueueFull(pub u64);

impl fmt::Display for QueueFull {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "queue full, rejected job {}", self.0)
    }
}

impl<'a> Queue<'a> {
    pub fn new(capacity: usize) -> Self {
        Self { heap: BinaryHeap::with_capacity(capacity), capacity }
    }

    pub fn push(&mut self, job: Job<'a>) -> Result<(), QueueFull> {
        if self.heap.len() >= self.capacity {
            return Err(QueueFull(job.id));
        }
        self.heap.push(job);
        Ok(())
    }

    /// Pops the next job, skipping ones older than `max_age`.
    pub fn pop_fresh(&mut self, max_age: Duration) -> Option<Job<'a>> {
        while let Some(job) = self.heap.pop() {
            if job.enqueued.elapsed() <= max_age {
                return Some(job);
            }
        }
        None
    }
}
