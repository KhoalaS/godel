<script setup lang="ts">
import { VueFlow } from '@vue-flow/core'
import CustomNode from './CustomNode.vue'
import { usePipelineStore } from '@/stores/pipeline'
import { Background } from '@vue-flow/background'
import CustomEdge from './CustomEdge.vue'
import { WindowSection } from 'vue-98'
import { useNodeDragAndDrop } from '@/composables/nodeDragAndDrop'

const store = usePipelineStore()

const { onDragOver, onDrop } = useNodeDragAndDrop(store)
</script>

<template>
  <div class="p-[1px] flex flex-1 min-h-0">
    <WindowSection
      class="flex-1 min-h-0 max-w-48 m-1 mb-0 ml-0 overflow-x-hidden overflow-y-auto p-2"
    >
      <div
        :key="category"
        class="mb-8"
        v-for="[category, nodes] in store.getCategorizedNodes.entries()"
      >
        <p>{{ category }}</p>
        <div
          class="node text-center p-3 mb-2 cursor-grab"
          draggable="true"
          @dragstart="(e) => e.dataTransfer?.setData('application/vueflow', node.type)"
          :key="node.type"
          v-for="node in nodes"
        >
          {{ node.name }}
        </div>
      </div>
    </WindowSection>
    <div class="flex-1 bg-white">
      <VueFlow @drop="onDrop" @dragover="onDragOver">
        <Background :size="1.6"></Background>
        <template #node-custom="nodeProps">
          <CustomNode v-bind="nodeProps" />
        </template>
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
