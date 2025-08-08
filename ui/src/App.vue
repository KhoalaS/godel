<script setup lang="ts">
import { ref } from 'vue'
import JobsTable from './components/JobsTable.vue'
import { useJobStore } from './stores/job'

const url = ref('')
const jobStore = useJobStore()

function download() {
  const parsedUrl = URL.parse(url.value)
  if (!parsedUrl) {
    console.warn('invalid url')
    return
  }

  jobStore.addJob(url.value, 'real-debrid')
  url.value = ''
}
</script>

<template>
  <label for="url" aria-label="Url">Url</label>
  <input v-model="url" id="url" type="text" />
  <button @click="download">Download</button>
  <Suspense>
    <!-- component with nested async dependencies -->
    <JobsTable />

    <!-- loading state via #fallback slot -->
    <template #fallback> Loading... </template>
  </Suspense>
</template>

<style scoped></style>
