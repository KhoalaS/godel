<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useVueFlow, VueFlow, type Edge, type Node } from '@vue-flow/core'
import CustomNode from './CustomNode.vue'
import { usePipelineStore } from '@/stores/pipeline'
import { Background } from '@vue-flow/background'

const pipelineStore = usePipelineStore()
const { addNodes, onConnect, addEdges, screenToFlowCoordinate } = useVueFlow()

// these components are only shown as examples of how to use a custom node or edge
// you can find many examples of how to create these custom components in the examples page of the docs

// these are our nodes
const nodes = ref<Node<{ label: string; hello?: string }>[]>([])

// these are our edges
const edges = ref<Edge[]>([])
onMounted(async () => {
  await pipelineStore.init()
})

onConnect((params) => {
  console.log(params)
  addEdges([params])
})

function onDragOver(event: DragEvent) {
  event.preventDefault()
  event.dataTransfer!.dropEffect = 'move'
}

function onDrop(event: DragEvent) {
  const type = event.dataTransfer?.getData('application/vueflow')
  const position = screenToFlowCoordinate({ x: event.clientX, y: event.clientY })

  const target = pipelineStore.registeredNodes.find((node) => node.type == type)
  if (target) {
    const id = crypto.randomUUID()
    addNodes({
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
  <div>
    {{ nodes }}
  </div>
  <div class="p-[1px]" style="height: 600px; width: 800px; display: flex">
    <div class="m-2">
      <div
        class="node w-32 text-center p-3"
        draggable="true"
        @dragstart="(e) => e.dataTransfer?.setData('application/vueflow', node.type)"
        style="cursor: grab"
        :key="node.type"
        v-for="node in pipelineStore.registeredNodes"
      >
        {{ node.name }}
      </div>
    </div>
    <div style="flex: 1; background-color: white">
      <VueFlow @drop="onDrop" @dragover="onDragOver" :nodes="nodes" :edges="edges">
        <Background :size="1.6"></Background>
        <!-- bind your custom node type to a component by using slots, slot names are always `node-<type>` -->
        <template #node-custom="nodeProps">
          <CustomNode v-bind="nodeProps" />
        </template>
        <template #node-test="nodeProps">
          <TestNode v-bind="nodeProps"></TestNode>
        </template>
        <!-- bind your custom edge type to a component by using slots, slot names are always `edge-<type>` -->
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
