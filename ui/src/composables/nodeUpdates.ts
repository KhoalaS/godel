import type { NodeIO, PipelineNode } from '@/models/Node'
import { FunctionRegistry } from '@/registries/InputHook'
import { usePipelineStore } from '@/stores/pipeline'
import { useNodeConnections } from '@vue-flow/core'

export function useNodeUpdates(nodeData: PipelineNode) {
  const data: PipelineNode = nodeData
  const store = usePipelineStore()
  const vueFlow = store.vueFlow

  const targetConnections = useNodeConnections({
    handleType: 'source',
  })

  function updateTargetNodes(
    inputs: {
      inputId: string
      newValue: string | number | boolean | undefined
    }[],
  ) {
    console.log('CustomNode:updateTargetNodes')

    const data = new Map<string, Record<string, NodeIO>>()

    for (const { inputId, newValue } of inputs) {
      if (newValue == undefined) {
        continue
      }

      const conns = targetConnections.value.filter((conn) => {
        return conn.sourceHandle == inputId
      })

      for (const conn of conns) {
        if (!conn.targetHandle) {
          return
        }

        const targetNode = vueFlow.findNode<PipelineNode>(conn.target)
        if (!targetNode || !targetNode.data.io) {
          continue
        }

        if (!data.has(conn.target)) {
          data.set(conn.target, {})
        }

        data.set(conn.target, {
          ...data.get(conn.target),
          [conn.targetHandle]: {
            ...targetNode.data.io[conn.targetHandle],
            value: newValue,
          },
        })
      }
    }

    for (const [nodeId, updates] of data.entries()) {
      const targetNode = vueFlow.findNode<PipelineNode>(nodeId)

      if (targetNode?.data.io) {
        vueFlow.updateNodeData<PipelineNode>(nodeId, {
          io: {
            ...targetNode.data.io,
            ...updates,
          },
        })
      }
    }
  }

  function onUpdate(value: string | number | boolean, io: NodeIO) {
    console.log('CustomNode:onUpdate', io.id)
    if (io.readOnly || data.io == null) {
      return
    }

    const hookUpdates: Record<string, NodeIO> = {}

    if (io.hooks) {
      for (const [hookId, functionId] of Object.entries(io.hooks)) {
        const func = FunctionRegistry.get(functionId)
        if (func == undefined) {
          continue
        }

        const hookValues: Record<string, unknown> = {}
        for (const _io of Object.values(data.io ?? {})) {
          if (_io.hookMapping?.[hookId] != undefined) {
            hookValues[_io.hookMapping[hookId]] = _io.id == io.id ? value : _io.value
          }
        }

        const newValue = func(hookValues)

        if (data.io?.[hookId].value == newValue) {
          continue
        }

        if (data.io) {
          hookUpdates[hookId] = {
            ...data.io?.[hookId],
            value: newValue,
          }
        }
      }
    }

    vueFlow.updateNodeData<PipelineNode>(data.id!, {
      io: {
        ...data.io,
        [io.id]: { ...io, value: value },
        ...hookUpdates,
      },
    })

    const inputs: {
      inputId: string
      newValue: string | number | boolean | undefined
    }[] = [{ inputId: io.id, newValue: value }]

    for (const [id, io] of Object.entries(hookUpdates)) {
      inputs.push({ inputId: id, newValue: io.value })
    }
    updateTargetNodes(inputs)
  }

  function onValueChange(io: NodeIO) {
    console.log('CustomNode:onValueChange', io.id)

    updateTargetNodes([
      {
        inputId: io.id,
        newValue: io.value,
      },
    ])
    if (io.value != undefined && io.hooks) {
      hook(io)
    }
  }

  function hook(io: NodeIO) {
    console.log('CustomNode:hook', io.value)

    const hookUpdates: Record<string, NodeIO> = {}

    for (const [hookId, functionId] of Object.entries(io.hooks ?? {})) {
      const func = FunctionRegistry.get(functionId)
      if (func == undefined) {
        continue
      }

      const hookValues: Record<string, unknown> = {}
      for (const io of Object.values(data.io ?? {})) {
        if (io.hookMapping?.[hookId] != undefined) {
          hookValues[io.hookMapping[hookId]] = io.value
        }
      }

      const newValue = func(hookValues)

      if (data.io?.[hookId].value == newValue) {
        continue
      }

      if (data.io) {
        hookUpdates[hookId] = {
          ...data.io?.[hookId],
          value: newValue,
        }
      }
    }

    if (data.io) {
      vueFlow.updateNodeData<PipelineNode>(data.id!, {
        io: {
          ...data.io,
          ...hookUpdates,
        },
      })
    }

    const inputs: {
      inputId: string
      newValue: string | number | boolean | undefined
    }[] = []
    for (const [id, io] of Object.entries(hookUpdates)) {
      inputs.push({ inputId: id, newValue: io.value })
    }
    updateTargetNodes(inputs)
  }

  return {
    updateTargetNodes,
    onUpdate,
    onValueChange,
  }
}
