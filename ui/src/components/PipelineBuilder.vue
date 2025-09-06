<script setup lang="ts">
import { onMounted } from 'vue'
import {
  useVueFlow,
  VueFlow,
  type Edge,
  type FlowExportObject,
  type GraphNode,
} from '@vue-flow/core'
import CustomNode from './CustomNode.vue'
import { usePipelineStore } from '@/stores/pipeline'
import { Background } from '@vue-flow/background'
import type { NodeIO, PipelineNode } from '@/types/Node'
import type { PipelineMessage } from '@/types/PipelineMessage'

const pipelineStore = usePipelineStore()
const {
  addNodes,
  onConnect,
  addEdges,
  screenToFlowCoordinate,
  onConnectStart,
  onConnectEnd,
  findNode,
  nodes,
  updateNodeData,
  toObject,
  fromObject,
  edges,
} = useVueFlow()

// these components are only shown as examples of how to use a custom node or edge
// you can find many examples of how to create these custom components in the examples page of the docs

// these are our edges
const messageCallback = (message: PipelineMessage) => {
  if (message.type == 'status') {
    switch (message.data.status) {
      case 'success':
        edges.value.forEach((e) => {
          if (e.source == message.nodeId) {
            e.animated = false
            e.label = undefined
          }
        })
        break
      case 'running':
        // set edge animations
        edges.value.forEach((e) => {
          if (e.source == message.nodeId) {
            e.animated = true
          }
        })
        break
      case 'failed':
        edges.value.forEach((e) => {
          if (e.source == message.nodeId) {
            e.animated = false
          }
        })
        break
    }
  } else if (message.type == 'progress') {
    edges.value.forEach((e) => {
      if (e.source == message.nodeId) {
        e.label = message.data.progress?.toFixed(2)
      }
    })
  } else if (message.type == 'error') {
    console.warn(message.data.error)
    edges.value.forEach((e) => {
      if (e.source == message.nodeId) {
        e.animated = false
      }
    })
  }
}
onMounted(async () => {
  await pipelineStore.init(messageCallback)
})

onConnect((params) => {
  console.log(params)
  const e: Edge = {
    ...params,
    id: crypto.randomUUID(),
    animated: false,
  }
  addEdges(e)

  const sourceNode = findNode<PipelineNode>(params.source)
  const targetNode = findNode<PipelineNode>(params.target)

  if (!sourceNode || !targetNode) {
    return
  }

  if (
    params.sourceHandle &&
    params.targetHandle &&
    sourceNode.data.io &&
    targetNode.data.io &&
    sourceNode.data.io[params.sourceHandle] &&
    targetNode.data.io[params.targetHandle]
  ) {
    const sourceData = sourceNode.data.io[params.sourceHandle].value
    if (sourceData != undefined) {
      updateNodeData<PipelineNode>(params.target, {
        io: {
          ...targetNode.data.io,
          [params.targetHandle]: {
            ...targetNode.data.io[params.targetHandle],
            value: sourceData,
          },
        },
      })
    }
  }
})

onConnectStart((params) => {
  // Find the source node
  const sourceNode = findNode(params.nodeId!)
  // Get the valueType of the source handle (adjust as needed for your node data structure)
  const sourceHandle = sourceNode?.data?.io?.[params.handleId!]

  nodes.value.forEach((node: GraphNode<PipelineNode>) => {
    if (node.id == params.nodeId) {
      return
    }

    const newIo: Record<string, NodeIO> = {}

    for (const [ioId, io] of Object.entries(node.data.io!)) {
      newIo[ioId] = {
        ...io,
        disabled: io.valueType != sourceHandle?.valueType || io.type == 'output',
      }
    }

    console.log(newIo)

    updateNodeData(node.id, {
      io: newIo,
    })
  })
})

onConnectEnd(() => {
  nodes.value.forEach((node: GraphNode<PipelineNode>) => {
    const newIo: Record<string, NodeIO> = {}

    for (const [ioId, io] of Object.entries(node.data.io!)) {
      newIo[ioId] = {
        ...io,
        disabled: false,
      }
    }

    console.log(newIo)

    updateNodeData(node.id, {
      io: newIo,
    })
  })
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

function startPipeline() {
  const graph = toObject()
  pipelineStore.startPipeline(graph)
}

function getPipelineObject() {
  return toObject()
}

function loadPipeline(obj: FlowExportObject) {
  return fromObject(obj)
}

defineExpose({
  saveGraph: startPipeline,
  getPipelineObject,
  loadPipeline,
})
</script>

<template>
  <div class="p-[1px]" style="height: 800px; width: 1600px; display: flex">
    <div class="m-2">
      <div
        :key="category"
        class="mb-8"
        v-for="[category, nodes] in pipelineStore.getCategorizedNodes.entries()"
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
