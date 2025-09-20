<script setup lang="ts">
import { VueFlow } from '@vue-flow/core'
import CustomNode from './CustomNode.vue'
import { usePipelineStore } from '@/stores/pipeline'
import { Background } from '@vue-flow/background'
import CustomEdge from './CustomEdge.vue'

const store = usePipelineStore()
const vueFlow = store.vueFlow

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
</script>

<template>
  <div class="p-[1px] flex flex-1">
    <div class="m-2 overflow-x-hidden overflow-y-auto pr-4">
      <div
        :key="category"
        class="mb-8"
        v-for="[category, nodes] in store.getCategorizedNodes.entries()"
      >
        <p>{{ category }}</p>
        <div
          class="node w-48 text-center p-3 mb-2 cursor-grab"
          draggable="true"
          @dragstart="(e) => e.dataTransfer?.setData('application/vueflow', node.type)"
          :key="node.type"
          v-for="node in nodes"
        >
          {{ node.name }}
        </div>
      </div>
    </div>
    <div class="flex-1 bg-white">
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
