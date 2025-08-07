<script setup lang="ts">
import JobsTable from './components/JobsTable.vue'
</script>

<template>
  <Suspense>
    <!-- component with nested async dependencies -->
    <JobsTable />

    <!-- loading state via #fallback slot -->
    <template #fallback> Loading... </template>
  </Suspense>
</template>

<style scoped></style>
