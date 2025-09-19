import { MessageLevel } from '@/messages/Message'
import z from 'zod'

export const PipelineMessageData = z.object({
  pipelineId: z.string(),
  status: z.enum(['started', 'done', 'failed']),
  error: z.string().optional(),
})

export const PipelineMessage = z.object({
  type: z.literal('pipelineUpdate'),
  data: PipelineMessageData,
  level: MessageLevel,
})

export type PipelineMessageData = z.infer<typeof PipelineMessageData>
export type PipelineMessage = z.infer<typeof PipelineMessage>
