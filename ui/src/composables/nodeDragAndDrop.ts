import type { PipelineStore } from '@/stores/pipeline'

export function useNodeDragAndDrop(store: PipelineStore) {
  const vueFlow = store.vueFlow

  function onDragOver(event: DragEvent) {
    event.preventDefault()
    event.dataTransfer!.dropEffect = 'move'
  }

  function onDrop(event: DragEvent) {
    const type = event.dataTransfer?.getData('application/vueflow')
    const position = vueFlow.screenToFlowCoordinate({ x: event.clientX, y: event.clientY })

    const target = store.registeredNodes.find((node) => node.type == type)
    if (target) {
      const id = crypto.randomUUID()
      vueFlow.addNodes({
        id: id,
        position: position,
        type: 'custom',
        data: {
          ...target,
          id,
        },
      })
    }
  }

  return { onDragOver, onDrop }
}
