import { useErrorService } from '@/composables/errorService'
import { PipelineMessageHandler } from '@/messages/PipelineMessageHandler'
import { HandleColors, NodeIO, PipelineNode } from '@/models/Node'
import { useVueFlow, type Edge, type FlowExportObject, type GraphNode } from '@vue-flow/core'
import { defineStore } from 'pinia'
import { computed, ref, type Ref } from 'vue'
import z from 'zod'

export const usePipelineStore = defineStore('pipeline', () => {
  const vueFlow = useVueFlow()
  const baseUrl = 'localhost:9095'
  const registeredNodes: Ref<PipelineNode[]> = ref([])
  const msgHandler = new PipelineMessageHandler(vueFlow)
  const errorService = useErrorService()

  async function init() {
    try {
      await initWs()
      const response = await fetch(`http://${baseUrl}/nodes`)
      if (response.status != 200) {
        return
      }

      const data = await response.json()
      registeredNodes.value = z.parse(z.array(PipelineNode), data)
    } catch (e: unknown) {
      errorService.handleError(e, 'error initializing pipeline store')
    }
  }

  async function initWs() {
    try {
      const socket = new WebSocket(`ws://${baseUrl}/updates/pipeline`)
      // Connection opened
      socket.addEventListener('open', () => {
        console.log('Connection opened')
      })

      socket.addEventListener('error', (event) => {
        console.warn('WebSocket error', event)
      })

      // Listen for messages
      socket.addEventListener('message', (event) => {
        try {
          const messageData = JSON.parse(event.data)

          msgHandler.onMessage(messageData)
        } catch (e: unknown) {
          errorService.handleError(e)
        }
      })
    } catch (e: unknown) {
      errorService.handleError(e, 'could not open websocket connection')
    }
  }

  async function startPipeline() {
    const graph: FlowExportObject = vueFlow.toObject()

    try {
      const response = await fetch(`http://${baseUrl}/pipeline/start`, {
        method: 'POST',
        body: JSON.stringify(graph),
      })
      if (response.status != 200) {
        return
      }
    } catch (e: unknown) {
      errorService.handleError(e)
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
    const sourceNode = vueFlow.findNode(params.nodeId!)
    // Get the valueType of the source handle (adjust as needed for your node data structure)
    const sourceHandle = sourceNode?.data?.io?.[params.handleId!]

    vueFlow.nodes.value.forEach((node: GraphNode<PipelineNode>) => {
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
