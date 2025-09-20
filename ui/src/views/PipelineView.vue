<script setup lang="ts">
import { WButton, WindowBody, WindowComponent } from 'vue-98'
import PipelineBuilder from '@/components/PipelineBuilder.vue'
import { usePipelineStore } from '@/stores/pipeline'

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

const store = usePipelineStore()

/**
 * Starts the current pipeline.
 */
function startPipeline() {
  store.startPipeline()
}

/**
 * Saves the current pipeline to a JSON file.
 */
function savePipeline() {
  const data = store.getPipelineObject()

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

  store.loadPipeline(graph)
}
</script>

<template>
  <WindowComponent
    class="h-full w-full"
    title="Pipeline Builder"
    :controls="['Minimize', 'Maximize']"
  >
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
