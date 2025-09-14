import type { IMessageHandler } from './IMessageHandler'
import type { VueFlowStore } from '@vue-flow/core'
import type { PipelineNode } from '@/models/Node'
import { PipelineMessage } from '@/models/messages/PipelineMessage'
import { NodeMessage } from '@/models/messages/NodeMessage'
import z from 'zod'

export class PipelineMessageHandler implements IMessageHandler<PipelineMessage | NodeMessage> {
  private vueFlow
  onMessage = (message: PipelineMessage | NodeMessage) => {
    if (message.type == 'nodeUpdate') {
      const _msg = z.parse(NodeMessage, message)
      this.nodeStatusHandler(_msg)
    } else if (message.type == 'pipelineUpdate') {
      const _msg = z.parse(PipelineMessage, message)
      this.pipelineStatusHandler(_msg)
    }
  }
  constructor(vueFlow: VueFlowStore) {
    this.vueFlow = vueFlow
  }

  private pipelineStatusHandler(message: PipelineMessage) {
    switch (message.data.status) {
      case 'done':
      case 'failed':
        this.vueFlow.setInteractive(true)
        break
      case 'started':
        this.vueFlow.setInteractive(false)
    }

    if (message.level == 'error') {
      console.error(message.data.error)
    }
  }

  private nodeStatusHandler(message: NodeMessage) {
    this.vueFlow.updateNodeData<PipelineNode>(message.data.nodeId, {
      status: message.data.status,
      progress: message.data.progress,
    })

    switch (message.data.status) {
      case 'running':
        // set edge animations
        this.vueFlow.edges.value.forEach((e) => {
          if (e.source == message.data.nodeId) {
            e.animated = true
          }
        })
        break
      case 'success':
      case 'failed':
        this.vueFlow.edges.value.forEach((e) => {
          if (e.source == message.data.nodeId) {
            e.animated = false
          }
        })
        break
    }

    if (message.level == 'error') {
      console.error(message.data.error)
      this.vueFlow.updateNodeData<PipelineNode>(message.data.nodeId, {
        error: message.data.error,
      })
      this.vueFlow.edges.value.forEach((e) => {
        if (e.source == message.data.nodeId) {
          e.animated = false
        }
      })
    }
  }
}
