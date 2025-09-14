<script setup lang="ts">
import { onMounted, useTemplateRef } from 'vue'
import { WButton, WindowBody, WindowComponent } from 'vue-98'
import PipelineBuilder from './components/PipelineBuilder.vue'

const pickerOpts = {
  types: [
    {
      description: 'JSON Files',
      accept: {
        'application/json': ['.json'],
      },
    },
  ],
  excludeAcceptAllOption: true,
  multiple: false,
}

/**
 * Starts the current pipeline.
 */
function startPipeline() {
  pipelineBuilder.value?.saveGraph()
}

/**
 * Saves the current pipeline to a JSON file.
 */
function savePipeline() {
  const data = pipelineBuilder.value?.getPipelineObject()
  if (!data) {
    return
  }

  const blobData = new Blob([JSON.stringify(data)], { type: 'application/json' })
  const url = URL.createObjectURL(blobData)
  const link = document.createElement('a')
  link.setAttribute('download', 'graph.json')
  link.href = url
  link.click()
}

/**
 * Loads a saved pipeline from a JSON file.
 */
async function loadPipeline() {
  const [fileHandle]: FileSystemFileHandle[] = await window.showOpenFilePicker(pickerOpts)
  if (fileHandle == undefined) {
    return
  }

  const text = await (await fileHandle.getFile()).text()
  const graph = JSON.parse(text)

  pipelineBuilder.value?.loadPipeline(graph)
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
            <WButton @click="savePipeline">Save Pipeline</WButton>
            <WButton @click="loadPipeline">Load Pipeline</WButton>
          </div>
        </template>
        <PipelineBuilder ref="pipelineBuilder"></PipelineBuilder>
      </WindowBody>
    </template>
  </WindowComponent>
</template>

<style scoped></style>
