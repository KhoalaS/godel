<script lang="ts" setup>
import type { NodeIO, PipelineNode } from '@/types/Node'
import { Handle, Position, useNodeConnections, useVueFlow, type NodeProps } from '@vue-flow/core'
import { WInput } from 'vue-98'
const props = defineProps<NodeProps<PipelineNode>>()

const { updateNodeData } = useVueFlow()

const _config: Record<string, unknown> = {}
if (props.data.outputs) {
  for (const key of Object.keys(props.data.outputs)) {
    const input = props.data.outputs[key]
    let value: unknown | undefined = input.value

    if (value == undefined) {
      switch (input.valueType) {
        case 'boolean':
          value = false
          break
        case 'number':
          value = 0
          break
        case 'string':
          value = ''
          break
        case 'directory':
          value = ''
      }
    }

    _config[input.id] = value
  }
}

const sourceConnections = useNodeConnections({
  // type target means all connections where *this* node is the target
  // that means we go backwards in the graph to find the source of the connection(s)
  handleType: 'target',
})

function onUpdate(value: string | number | boolean, type: 'input' | 'output', io: NodeIO) {
  if (type == 'output' && !io.readOnly && props.data.outputs) {
    updateNodeData<PipelineNode>(props.id, {
      outputs: {
        ...props.data.outputs,
        [io.id]: { ...props.data.outputs?.[io.id], value: value },
      },
    })
  } else if (type === 'input' && !io.readOnly && props.data.inputs) {
    updateNodeData<PipelineNode>(props.id, {
      inputs: {
        ...props.data.inputs,
        [io.id]: { ...props.data.inputs?.[io.id], value: value },
      },
    })
  }
}
</script>

<template>
  <div class="node w-64 text-center p-3">
    <div class="input-wrapper" :key="input.id" v-for="input in data.inputs">
      <Handle
        class="handle-input"
        :key="input.id"
        type="target"
        :id="input.id"
        :position="Position.Left"
        :connectable-start="false"
        :connectable-end="true"
      />
      <label>{{ input.label }}</label>
      <WInput @update="(v) => onUpdate(v, 'input', input)" :type="input.valueType" />
    </div>
    <div class="input-wrapper" :key="output.id" v-for="output in data.outputs">
      <Handle
        class="handle"
        :key="output.id"
        type="source"
        :id="output.id"
        :position="Position.Right"
        :connectable-start="true"
        :connectable-end="false"
      />
      <label>{{ output.label }}</label>
      <WInput @update="(v) => onUpdate(v, 'output', output)" :type="output.valueType" />
    </div>
  </div>
</template>

<style scoped>
.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.handle-input {
  left: -12px;
  position: absolute;
}

.node {
  box-shadow:
    inset -1px -1px black,
    inset 1px 1px white,
    inset -2px -2px var(--border-gray);
  background-color: var(--main-bg-color);
}
</style>
