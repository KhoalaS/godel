<script setup lang="ts">
import { useJobStore } from '@/stores/job'
import JobComponent from './JobComponent.vue'

const jobStore = useJobStore()
</script>
<template>
  <div>
    <JobComponent v-for="job in jobStore.jobs" :key="job.id" :job="job"></JobComponent>
  </div>
</template>
