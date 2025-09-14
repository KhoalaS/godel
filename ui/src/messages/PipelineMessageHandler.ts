import type { PipelineMessage } from '@/models/PipelineMessage'
import type { IMessageHandler } from './IMessageHandler'
import type { VueFlowStore } from '@vue-flow/core'
import type { PipelineNode } from '@/models/Node'

export class PipelineMessageHandler implements IMessageHandler<PipelineMessage> {
  private vueFlow
  onMessage = (message: PipelineMessage) => {
    if (message.type == 'status') {
      this.vueFlow.updateNodeData<PipelineNode>(message.nodeId, {
        status: message.data.status,
        progress: undefined,
      })
      switch (message.data.status) {
        case 'success':
          this.vueFlow.edges.value.forEach((e) => {
            if (e.source == message.nodeId) {
              e.animated = false
            }
          })
          break
        case 'running':
          // set edge animations
          this.vueFlow.edges.value.forEach((e) => {
            if (e.source == message.nodeId) {
              e.animated = true
            }
          })
          break
        case 'failed':
          this.vueFlow.edges.value.forEach((e) => {
            if (e.source == message.nodeId) {
              e.animated = false
            }
          })
          break
      }
    } else if (message.type == 'progress') {
      this.vueFlow.updateNodeData<PipelineNode>(message.nodeId, {
        progress: message.data.progress,
        status: message.data.status,
      })
    } else if (message.type == 'error') {
      console.warn(message.data.error)

      this.vueFlow.updateNodeData<PipelineNode>(message.nodeId, {
        status: message.data.status,
        error: message.data.error,
        progress: undefined,
      })
      this.vueFlow.edges.value.forEach((e) => {
        if (e.source == message.nodeId) {
          e.animated = false
        }
      })
    }
  }
  constructor(vueFlow: VueFlowStore) {
    this.vueFlow = vueFlow
  }
}
