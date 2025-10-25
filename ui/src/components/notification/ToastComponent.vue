<script setup lang="ts">
import type { NotificationConfig } from '@/services/NotificationService'
import { onMounted } from 'vue'
import { WindowBody, WindowComponent } from 'vue-98'

const SECOND_IN_MILLISECONDS = 1000

const props = defineProps<NotificationConfig>()

const emit = defineEmits<{
  close: []
}>()

let timeoutHandle: number | undefined

const onControlClick = () => {
  clearTimeout(timeoutHandle)
  emit('close')
}

onMounted(() => {
  if (props.durationInSeconds) {
    timeoutHandle = setTimeout(() => {
      emit('close')
    }, props.durationInSeconds * SECOND_IN_MILLISECONDS)
  }
})
</script>
<template>
  <WindowComponent
    :title="props.title"
    :controls="['Close']"
    @click:control="onControlClick"
    class="w-48 h-24 z-50 absolute top-0 left-0 bg-amber-600"
  >
    <template #body>
      <WindowBody class="p-2">
        <div>{{ props.message }}</div>
      </WindowBody>
    </template>
  </WindowComponent>
</template>
