<script setup lang="ts">
const title: string = `stateful ${"highlighting"}`
</script>

<template>
  <article
    class="card"
    :aria-label="title"
  >{{ title }}</article>
</template>

<style scoped>
.card { color: #58a6ff; }
</style>
