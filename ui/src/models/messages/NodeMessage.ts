import { MessageLevel } from '@/messages/Message'
import z from 'zod'

export const NodeMessageData = z.object({
  pipelineId: z.string(),
  nodeId: z.string(),
  error: z.string().optional(),
  progress: z.float64().optional(),
  status: z.enum(['pending', 'running', 'success', 'failed']),
})

export const NodeMessage = z.object({
  type: z.literal('nodeUpdate'),
  data: NodeMessageData,
  level: MessageLevel,
})

export type NodeMessageData = z.infer<typeof NodeMessageData>
export type NodeMessage = z.infer<typeof NodeMessage>
