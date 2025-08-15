<script setup lang="ts">
import { useJobStore } from '@/stores/job'
import type { DownloadJob } from '@/types/DownloadJob'
import { computed, type DeepReadonly } from 'vue'
import WButton from 'vue-98'

const props = defineProps<{
  job: DeepReadonly<DownloadJob>
}>()

const jobStore = useJobStore()

const pauseButtonState = computed(() => {
  return props.job.status != 'downloading' && props.job.status != 'paused'
})

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

const speed = computed(() => {
  switch (props.job.status) {
    case 'downloading':
      if (!props.job.speed) {
        return '-'
      } else {
        return `${(props.job.speed / 1024 / 1024).toFixed(2)} MB/s`
      }
    default:
      return ''
  }
})

const eta = computed(() => {
  switch (props.job.status) {
    case 'downloading':
      return props.job.eta?.toFixed(2) ?? '-'
    default:
      return '-'
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
  <WButton :disabled="pauseButtonState" @click="handlePause">{{ pauseLabel }}</WButton>
  <span>Speed: {{ speed }}</span>
  <span>ETA: {{ eta }} Seconds</span>
</template>
