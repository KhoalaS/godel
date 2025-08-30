<script lang="ts" setup>
import type { NodeIO, PipelineNode } from '@/types/Node'
import {
  Handle,
  Position,
  useNodeConnections,
  useNodesData,
  useVueFlow,
  type NodeProps,
} from '@vue-flow/core'
import { WInput } from 'vue-98'
const props = defineProps<NodeProps<PipelineNode>>()

const { updateNodeData } = useVueFlow()

const sourceConnections = useNodeConnections({
  // type target means all connections where *this* node is the target
  // that means we go backwards in the graph to find the source of the connection(s)
  handleType: 'target',
})

const sourceData = useNodesData(() =>
  sourceConnections.value.map((connection) => connection.source),
)

function onUpdate(value: string | number | boolean, io: NodeIO) {
  if (!io.readOnly && props.data.io) {
    updateNodeData<PipelineNode>(props.id, {
      io: {
        ...props.data.io,
        [io.id]: { ...props.data.io?.[io.id], value: value },
      },
    })
  }
}
</script>

<template>
  <div class="node w-64 text-center p-3">
    <div class="input-wrapper" :key="input.id" v-for="input in data.io">
      <Handle
        v-if="input.type == 'input' || input.type == 'passthrough'"
        class="handle-input"
        :key="input.id"
        type="target"
        :id="input.id"
        :position="Position.Left"
        :connectable-start="false"
        :connectable-end="true"
      />
      <Handle
        v-if="input.type == 'output' || input.type == 'passthrough'"
        class="handle-output"
        :key="input.id"
        type="source"
        :id="input.id"
        :position="Position.Right"
        :connectable-start="true"
        :connectable-end="false"
      />
      <label>{{ input.label }}</label>
      <WInput @update="(v) => onUpdate(v, input)" :type="input.valueType" />
    </div>
  </div>
  {{ sourceConnections }}
  {{ sourceData }}
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

.handle-output {
  right: -12px;
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
