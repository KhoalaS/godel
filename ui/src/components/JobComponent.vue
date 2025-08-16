<script setup lang="ts">
import { useJobStore } from '@/stores/job'
import type { DownloadJob } from '@/types/DownloadJob'
import { computed, type DeepReadonly } from 'vue'
import { ProgressbarComponent, WButton } from 'vue-98'

const props = defineProps<{
  job: DeepReadonly<DownloadJob>
}>()

const jobStore = useJobStore()
const progress = computed(() => {
  return (props.job.bytesDownloaded ?? 0) / (props.job.size ?? 100)
})

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
      return '-'
  }
})

const eta = computed(() => {
  switch (props.job.status) {
    case 'downloading':
      return props.job.eta ? props.job.eta?.toFixed(2) + ' Seconds' : '-'
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

async function copyLink() {
  try {
    await navigator.clipboard.writeText(props.job.url)
    console.log('Link copied to clipboard')
  } catch (err) {
    console.error('Failed to copy link: ', err)
  }
}
</script>
<template>
  <div>
    <span>{{ job.url }}</span>
    <WButton @click="copyLink" style="display: inline-block"> Copy </WButton>
  </div>
  <ProgressbarComponent
    progress-label="Progress"
    :model-value="progress"
    :max="1"
  ></ProgressbarComponent>
  <WButton v-if="!pauseButtonState" @click="handlePause">{{ pauseLabel }}</WButton>
  <span>Speed: {{ speed }}</span>
  <span>ETA: {{ eta }} </span>
</template>
