<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useVueFlow, VueFlow, type Edge, type Node } from '@vue-flow/core'
import LimiterNode from './LimiterNode.vue'
import { usePipelineStore } from '@/stores/pipeline'
import { Background } from '@vue-flow/background'

const pipelineStore = usePipelineStore()
const { addNodes, onConnect, addEdges, screenToFlowCoordinate } = useVueFlow()

// these components are only shown as examples of how to use a custom node or edge
// you can find many examples of how to create these custom components in the examples page of the docs

// these are our nodes
const nodes = ref<Node<{ label: string; hello?: string }>[]>([
  // an input node, specified by using `type: 'input'`
  {
    id: '1',
    position: { x: 250, y: 5 },
    // all nodes can have a data object containing any data you want to pass to the node
    // a label can property can be used for default nodes
    data: { label: 'Node 1' },
    connectable: true,
  },

  // default node, you can omit `type: 'default'` as it's the fallback type
  {
    id: '2',
    position: { x: 100, y: 100 },
    data: { label: 'Node 2' },
  },

  // An output node, specified by using `type: 'output'`
  {
    id: '3',
    type: 'output',
    position: { x: 400, y: 200 },
    data: { label: 'Node 3' },
  },

  // this is a custom node
  // we set it by using a custom type name we choose, in this example `special`
  // the name can be freely chosen, there are no restrictions as long as it's a string
  {
    id: '4',
    position: { x: 400, y: 200 },
    data: {
      label: 'Node 4',
      hello: 'world',
    },
    connectable: true,
  },
  {
    id: '5',
    position: { x: 500, y: 200 },
    data: {
      label: 'Node 5',
    },
    connectable: true,
  },
])

// these are our edges
const edges = ref<Edge[]>([
  // default bezier edge
  // consists of an edge id, source node id and target node id
  {
    id: 'e1->2',
    source: '1',
    target: '2',
  },

  // set `animated: true` to create an animated edge path
  {
    id: 'e2->3',
    source: '2',
    target: '3',
    animated: true,
  },

  // a custom edge, specified by using a custom type name
  // we choose `type: 'special'` for this example
  {
    id: 'e3->4',
    type: 'step',
    source: '3',
    target: '4',

    // all edges can have a data object containing any data you want to pass to the edge
    data: {
      hello: 'world',
    },
  },
])
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
      type: type,
      data: {
        id,
        ...target,
      },
    })
  }
}
</script>

<template>
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
        <template #node-limiter="nodeProps">
          <LimiterNode v-bind="nodeProps" />
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
