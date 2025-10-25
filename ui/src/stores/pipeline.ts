import { PipelineMessageHandler } from '@/messages/PipelineMessageHandler'
import { HandleColors, NodeIO, PipelineNode } from '@/models/Node'
import { useVueFlow, type Edge, type FlowExportObject, type GraphNode } from '@vue-flow/core'
import { defineStore } from 'pinia'
import { computed, ref, type Ref } from 'vue'
import z from 'zod'

export const usePipelineStore = defineStore('pipeline', () => {
  const vueFlow = useVueFlow()
  const registeredNodes: Ref<PipelineNode[]> = ref([])
  const msgHandler = new PipelineMessageHandler(vueFlow)
  let initialized = false

  async function init() {
    if (initialized) {
      console.warn('tried re-initializing pipeline store')
      return
    }

    await initWs()
    const response = await fetch('/nodes', {
      headers: {
        accept: 'application/json',
      },
    })
    if (response.status != 200) {
      throw new Error(await response.text())
    }

    const data = await response.json()
    registeredNodes.value = z.parse(z.array(PipelineNode), data)

    initialized = true
  }

  async function initWs() {
    return new Promise((resolve, reject) => {
      const socket = new WebSocket(`ws://${window.location.host}/updates/pipeline`)
      // Connection opened
      socket.addEventListener('open', () => {
        console.log('Connection opened')
        resolve(socket)
      })

      socket.addEventListener('error', () => {
        console.warn('event')
        reject(new Error('WebSocket connection failed'))
      })

      // Listen for messages
      socket.addEventListener('message', (event) => {
        const messageData = JSON.parse(event.data)

        msgHandler.onMessage(messageData)
      })
    })
  }

  async function startPipeline() {
    const graph: FlowExportObject = vueFlow.toObject()

    const response = await fetch('/pipeline/start', {
      method: 'POST',
      body: JSON.stringify(graph),
    })
    if (response.status != 200) {
      return
    }
  }

  function getPipelineObject() {
    return vueFlow.toObject()
  }

  function loadPipeline(obj: FlowExportObject) {
    return vueFlow.fromObject(obj)
  }

  const getCategorizedNodes = computed(() => {
    const acc = new Map<string, PipelineNode[]>()
    acc.set('other', [])

    registeredNodes.value.forEach((node) => {
      if (node.category) {
        if (acc.has(node.category)) {
          acc.get(node.category)!.push(node)
        } else {
          acc.set(node.category, [node])
        }
      } else {
        acc.get('other')!.push(node)
      }
    })

    return acc
  })

  vueFlow.onConnect((params) => {
    const sourceNode = vueFlow.findNode<PipelineNode>(params.source)
    const targetNode = vueFlow.findNode<PipelineNode>(params.target)

    const targetHandleType = targetNode?.data.io?.[params.targetHandle ?? ''].valueType

    const e: Edge = {
      ...params,
      id: crypto.randomUUID(),
      animated: false,
      type: 'custom',
      style: {
        stroke: targetHandleType ? HandleColors[targetHandleType] : undefined,
      },
    }
    vueFlow.addEdges(e)

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
        vueFlow.updateNodeData<PipelineNode>(params.target, {
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

  /**
   * vue-flow connectStart-handler. Disables handles that do not have the same handle type.
   */
  vueFlow.onConnectStart((params) => {
    // Find the source node
    const sourceNode = vueFlow.findNode<PipelineNode>(params.nodeId!)
    // Get the valueType of the source handle (adjust as needed for your node data structure)
    const sourceHandle = sourceNode?.data?.io?.[params.handleId!]

    vueFlow.nodes.value.forEach((node: GraphNode<PipelineNode>) => {
      const newIo: Record<string, NodeIO> = {}
      const isCurrentNode = node.id == sourceNode?.id

      for (const [ioId, io] of Object.entries(node.data.io!)) {
        const isSourceHandle = isCurrentNode && ioId == params.handleId
        if (isSourceHandle) {
          newIo[ioId] = io
          continue
        }

        newIo[ioId] = {
          ...io,
          disabled: io.valueType != sourceHandle?.valueType || io.type == 'output' || isCurrentNode,
        }
      }

      vueFlow.updateNodeData(node.id, {
        io: newIo,
      })
    })
  })

  /**
   * vue-flow connectEnd-handler. Adds a new edge.
   */
  vueFlow.onConnectEnd(() => {
    vueFlow.nodes.value.forEach((node: GraphNode<PipelineNode>) => {
      const newIo: Record<string, NodeIO> = {}

      for (const [ioId, io] of Object.entries(node.data.io!)) {
        newIo[ioId] = {
          ...io,
          disabled: false,
        }
      }

      vueFlow.updateNodeData(node.id, {
        io: newIo,
      })
    })
  })

  return {
    getCategorizedNodes,
    getPipelineObject,
    init,
    loadPipeline,
    registeredNodes,
    startPipeline,
    vueFlow,
  }
})

export type PipelineStore = ReturnType<typeof usePipelineStore>
