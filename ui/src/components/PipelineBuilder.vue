<script setup lang="ts">
import { onMounted } from 'vue'
import { VueFlow, type FlowExportObject } from '@vue-flow/core'
import CustomNode from './CustomNode.vue'
import { usePipelineStore } from '@/stores/pipeline'
import { Background } from '@vue-flow/background'
import { type PipelineNode } from '@/models/Node'
import type { PipelineMessage } from '@/models/PipelineMessage'
import CustomEdge from './CustomEdge.vue'

const store = usePipelineStore()
const vueFlow = store.vueFlow

// these components are only shown as examples of how to use a custom node or edge
// you can find many examples of how to create these custom components in the examples page of the docs

// these are our edges
const messageCallback = (message: PipelineMessage) => {
  if (message.type == 'status') {
    vueFlow.updateNodeData<PipelineNode>(message.nodeId, {
      status: message.data.status,
      progress: undefined,
    })
    switch (message.data.status) {
      case 'success':
        vueFlow.edges.forEach((e) => {
          if (e.source == message.nodeId) {
            e.animated = false
          }
        })
        break
      case 'running':
        // set edge animations
        vueFlow.edges.forEach((e) => {
          if (e.source == message.nodeId) {
            e.animated = true
          }
        })
        break
      case 'failed':
        vueFlow.edges.forEach((e) => {
          if (e.source == message.nodeId) {
            e.animated = false
          }
        })
        break
    }
  } else if (message.type == 'progress') {
    vueFlow.updateNodeData<PipelineNode>(message.nodeId, {
      progress: message.data.progress,
      status: message.data.status,
    })
  } else if (message.type == 'error') {
    console.warn(message.data.error)
    vueFlow.updateNodeData<PipelineNode>(message.nodeId, {
      status: message.data.status,
      error: message.data.error,
      progress: undefined,
    })
    vueFlow.edges.forEach((e) => {
      if (e.source == message.nodeId) {
        e.animated = false
      }
    })
  }
}
onMounted(async () => {
  await store.init(messageCallback)
})

function onDragOver(event: DragEvent) {
  event.preventDefault()
  event.dataTransfer!.dropEffect = 'move'
}

function onDrop(event: DragEvent) {
  const type = event.dataTransfer?.getData('application/vueflow')
  const position = vueFlow.screenToFlowCoordinate({ x: event.clientX, y: event.clientY })

  const target = store.registeredNodes.find((node) => node.type == type)
  if (target) {
    const id = crypto.randomUUID()
    vueFlow.addNodes({
      id: id,
      position: position,
      type: 'custom',
      data: {
        ...target,
        id,
      },
    })
  }
}

function startPipeline() {
  const graph = vueFlow.toObject()
  store.startPipeline(graph)
}

function getPipelineObject() {
  return vueFlow.toObject()
}

function loadPipeline(obj: FlowExportObject) {
  return vueFlow.fromObject(obj)
}

defineExpose({
  saveGraph: startPipeline,
  getPipelineObject,
  loadPipeline,
})
</script>

<template>
  <div class="p-[1px]" style="height: 800px; width: 1600px; display: flex">
    <div class="m-2 overflow-x-hidden overflow-y-auto pr-4">
      <div
        :key="category"
        class="mb-8"
        v-for="[category, nodes] in store.getCategorizedNodes.entries()"
      >
        <p>{{ category }}</p>
        <div
          class="node w-32 text-center p-3 mb-2"
          draggable="true"
          @dragstart="(e) => e.dataTransfer?.setData('application/vueflow', node.type)"
          style="cursor: grab"
          :key="node.type"
          v-for="node in nodes"
        >
          {{ node.name }}
        </div>
      </div>
    </div>
    <div style="flex: 1; background-color: white">
      <VueFlow @drop="onDrop" @dragover="onDragOver">
        <Background :size="1.6"></Background>
        <!-- bind your custom node type to a component by using slots, slot names are always `node-<type>` -->
        <template #node-custom="nodeProps">
          <CustomNode v-bind="nodeProps" />
        </template>
        <!-- bind your custom edge type to a component by using slots, slot names are always `edge-<type>` -->
        <template #edge-custom="customEdgeProps">
          <CustomEdge v-bind="customEdgeProps" />
        </template>
      </VueFlow>
    </div>
  </div>
</template>

<style scoped>
.node {
  box-shadow:
    inset -1px -1px black,
    inset 1px 1px white,
    inset -2px -2px var(--border-gray);
  background-color: var(--main-bg-color);
}
</style>
