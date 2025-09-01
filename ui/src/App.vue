<script setup lang="ts">
import { onMounted, useTemplateRef } from 'vue'
import { WButton, WindowBody, WindowComponent } from 'vue-98'
import PipelineBuilder from './components/PipelineBuilder.vue'

const pipelineBuilder = useTemplateRef('pipelineBuilder')

function startPipeline() {
  pipelineBuilder.value?.saveGraph()
}

onMounted(async () => {})
</script>

<template>
  <WindowComponent title="Pipeline Builder" :controls="['Minimize', 'Maximize']">
    <template #body>
      <WindowBody>
        <template #toolbars>
          <div class="flex gap-1 items-center py-1 px-[2px] justify-start">
            <WButton @click="startPipeline">Start</WButton>
          </div>
        </template>
        <PipelineBuilder ref="pipelineBuilder"></PipelineBuilder>
      </WindowBody>
    </template>
  </WindowComponent>
</template>

<style scoped></style>
