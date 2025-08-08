<script setup lang="ts">
import { useJobStore } from '@/stores/job'
import type { DownloadJob } from '@/types/DownloadJob'
import { computed, type DeepReadonly } from 'vue'

const props = defineProps<{
  job: DeepReadonly<DownloadJob>
}>()

const jobStore = useJobStore()

const pauseLabel = computed(() => {
  switch (props.job.status) {
    case 'paused':
      return 'Resume'
    case 'downloading':
      return 'Pause'
    default:
      return ''
  }
})

function handlePause() {
  switch (props.job.status) {
    case 'paused':
    case 'downloading':
      jobStore.pauseJob(props.job.id)
      return
    default:
      return
  }
}
</script>
<template>
  <div>{{ job.url }}</div>
  <progress :value="job.bytesDownloaded" :max="job.size"></progress>
  <button @click="handlePause">{{ pauseLabel }}</button>
</template>
