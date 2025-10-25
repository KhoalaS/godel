<script setup lang="ts">
import { RouterView } from 'vue-router'
import { usePipelineStore } from './stores/pipeline'
import { onMounted } from 'vue'

const store = usePipelineStore()
onMounted(async () => {
  await store.init()
})
</script>

<template>
  <RouterView />
</template>

<style scoped></style>
