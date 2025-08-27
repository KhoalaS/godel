<script lang="ts" setup>
import type { PipelineNode } from '@/types/Node'
import { Handle, Position, type NodeProps } from '@vue-flow/core'
import { reactive } from 'vue'
import { WInput } from 'vue-98'
const props = defineProps<NodeProps<PipelineNode>>()

const _config: Record<string, string | number | boolean> = {}
if (props.data.outputs) {
  for (const key of Object.keys(props.data.outputs)) {
    const input = props.data.outputs[key]
    let value: string | number | boolean

    switch (input.type) {
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

    _config[input.id] = value
  }
}

const config = reactive(_config)
</script>

<template>
  <div class="node w-32 text-center p-3">
    {{ config }}
    <Handle
      :key="value.id"
      v-for="value in data.inputs"
      type="target"
      :id="value.id"
      :position="Position.Left"
      :connectable-start="false"
      :connectable-end="true"
    />
    <Handle
      :key="value.id"
      v-for="value in data.outputs"
      type="source"
      :id="value.id"
      :position="Position.Right"
      :connectable-start="true"
      :connectable-end="false"
    />
    <div class="text-left" :key="input.id" v-for="input in data.inputs">
      <label>{{ input.label }}</label>
      <WInput v-model="config[input.id]" :type="input.type" />
    </div>
    <div class="text-left" :key="output.id" v-for="output in data.outputs">
      <label>{{ output.label }}</label>
      <WInput v-model="config[output.id]" :type="output.type" />
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
