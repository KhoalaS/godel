import z from 'zod'

export const PipelineMessage = z.object({
  pipelineId: z.string(),
  nodeId: z.string(),
  nodeType: z.string(),
  type: z.enum(['error', 'progress', 'status']),
  data: z.object({
    error: z.string().optional(),
    progress: z.float64().optional(),
    status: z.enum(['pending', 'running', 'success', 'failed']),
  }),
})

export type PipelineMessage = z.infer<typeof PipelineMessage>
