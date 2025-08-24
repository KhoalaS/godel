<script lang="ts" setup>
import type { PipelineNode } from '@/types/Node'
import { Handle, Position, type NodeProps } from '@vue-flow/core'
const props = defineProps<NodeProps<PipelineNode>>()
</script>

<template>
  <div class="downloader-node">
    <Handle type="source" id="output" :position="Position.Bottom" />
    <div>
      {{ props.data.name }}
    </div>
  </div>
</template>

<style>
.downloader-node {
  min-width: 100px;
  gap: 4px;
  padding: 8px;
  background: white;
  border: 1px solid blue;
  border-radius: 4px;
}
</style>
