import { PipelineNode } from '@/types/Node'
import type { FlowExportObject } from '@vue-flow/core'
import { defineStore } from 'pinia'
import { ref, type Ref } from 'vue'
import z from 'zod'

export const usePipelineStore = defineStore('pipeline', () => {
  const baseUrl = 'localhost:9095'
  const registeredNodes: Ref<PipelineNode[]> = ref([])
  const pipelines = new Map<string, unknown>()

  async function init() {
    try {
      const response = await fetch(`http://${baseUrl}/nodes`)
      if (response.status != 200) {
        return
      }

      const data = await response.json()
      registeredNodes.value = z.parse(z.array(PipelineNode), data)
    } catch (e: unknown) {
      console.log(e)
    }
  }

  async function startPipeline(graph: FlowExportObject) {
    try {
      const response = await fetch(`http://${baseUrl}/pipeline/start`, {
        method: 'POST',
        body: JSON.stringify(graph),
      })
      if (response.status != 200) {
        return
      }
    } catch (e: unknown) {
      console.log(e)
    }
  }

  return {
    init,
    registeredNodes,
    startPipeline,
  }
})
