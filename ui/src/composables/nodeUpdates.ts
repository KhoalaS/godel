import type { NodeIO, PipelineNode } from '@/models/Node'
import type { PipelineStore } from '@/stores/pipeline'
import { useNodeConnections, type NodeProps } from '@vue-flow/core'
import { useHookFunctionService } from './hookFunctionService'

export function useNodeUpdates(props: NodeProps<PipelineNode>, pipelineStore: PipelineStore) {
  const vueFlow = pipelineStore.vueFlow
  const hookFunctionService = useHookFunctionService()

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

  /**
   * Call this inside the node component whenever it's model value is updated.
   *
   * @param io The NodeIO object of the node.
   */
  function onUpdate(value: string | number | boolean, io: NodeIO) {
    console.log('CustomNode:onUpdate', io.id)
    if (io.readOnly || props.data.io == null) {
      return
    }

    const hookUpdates: Record<string, NodeIO> = {}

    if (io.hooks) {
      for (const [hookId, functionId] of Object.entries(io.hooks)) {
        const func = hookFunctionService.getFunction(functionId)
        if (func == undefined) {
          continue
        }

        const hookValues: Record<string, unknown> = {}
        for (const _io of Object.values(props.data.io ?? {})) {
          if (_io.hookMapping?.[hookId] != undefined) {
            hookValues[_io.hookMapping[hookId]] = _io.id == io.id ? value : _io.value
          }
        }

        const newValue = func(hookValues)

        if (props.data.io?.[hookId].value == newValue) {
          continue
        }

        if (props.data.io) {
          hookUpdates[hookId] = {
            ...props.data.io?.[hookId],
            value: newValue,
          }
        }
      }
    }

    vueFlow.updateNodeData<PipelineNode>(props.data.id!, {
      io: {
        ...props.data.io,
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

  /**
   * Call this inside the node component whenever it's value is updated with
   * the VueFlow API.
   *
   * @param io The NodeIO object of the node.
   */
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
      const func = hookFunctionService.getFunction(functionId)
      if (func == undefined) {
        continue
      }

      const hookValues: Record<string, unknown> = {}
      for (const io of Object.values(props.data.io ?? {})) {
        if (io.hookMapping?.[hookId] != undefined) {
          hookValues[io.hookMapping[hookId]] = io.value
        }
      }

      const newValue = func(hookValues)

      if (props.data.io?.[hookId].value == newValue) {
        continue
      }

      if (props.data.io) {
        hookUpdates[hookId] = {
          ...props.data.io?.[hookId],
          value: newValue,
        }
      }
    }

    if (props.data.io) {
      vueFlow.updateNodeData<PipelineNode>(props.data.id!, {
        io: {
          ...props.data.io,
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
