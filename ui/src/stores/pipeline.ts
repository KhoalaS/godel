import { PipelineNode } from '@/types/Node'
import type { FlowExportObject } from '@vue-flow/core'
import { defineStore } from 'pinia'
import { computed, ref, type Ref } from 'vue'
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

  return {
    init,
    registeredNodes,
    startPipeline,
    getCategorizedNodes,
  }
})
