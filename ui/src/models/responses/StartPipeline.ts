import z from 'zod'

export const StartPipelineResponse = z.object({
  pipelineId: z.string(),
})

export type StartPipelineResponse = z.infer<typeof StartPipelineResponse>
