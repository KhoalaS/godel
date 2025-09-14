<script setup lang="ts">
import { usePipelineStore } from '@/stores/pipeline'
import { BaseEdge, EdgeLabelRenderer, getBezierPath, type EdgeProps } from '@vue-flow/core'
import { computed, ref } from 'vue'
import { WindowButton } from 'vue-98'

const props = defineProps<EdgeProps>()
const { vueFlow } = usePipelineStore()

const path = computed(() => getBezierPath(props))
const hovered = ref(false)

function onRemoveEdge(id: string) {
  vueFlow.removeEdges(id)
}
</script>
<template>
  <g @mouseenter="hovered = true" @mouseleave="hovered = false">
    <BaseEdge
      :id="id"
      :style="{
        'stroke-width': hovered ? '2' : undefined,
        ...style,
      }"
      :path="path[0]"
    />
    <EdgeLabelRenderer>
      <div
        @mouseenter="hovered = true"
        @mouseleave="hovered = false"
        :style="{
          transform: `translate(-50%, -50%) translate(${path[1]}px,${path[2]}px)`,
        }"
        class="absolute nodrag nopan"
      >
        <WindowButton
          style="pointer-events: all"
          class="duration-100"
          v-show="hovered"
          @click="onRemoveEdge(id)"
          type="Close"
        ></WindowButton>
      </div>
    </EdgeLabelRenderer>
  </g>
</template>

<style scoped></style>
