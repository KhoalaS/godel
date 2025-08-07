<script setup lang="ts">
import type { DownloadJob } from '@/types/DownloadJob'
import type { DeepReadonly } from 'vue'

defineProps<{
  job: DeepReadonly<DownloadJob>
}>()
</script>
<template>
  <div>{{ job.url }}</div>
  <progress :value="job.bytesDownloaded" :max="job.size"></progress>
</template>
