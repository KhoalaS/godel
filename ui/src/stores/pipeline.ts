import { PipelineNode } from '@/types/Node'
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

  async function startPipeline(pipeline: unknown) {
    throw new Error('Unimplemented')
  }

  return {
    init,
    registeredNodes,
    startPipeline,
  }
})
