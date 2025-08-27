import z from 'zod'

export const NodeInput = z.object({
  id: z.string(),
  type: z.enum(['string', 'number', 'boolean']),
  label: z.string(),
  required: z.boolean(),
  options: z.array(z.string()).optional(),
})

export const PipelineNode = z.object({
  id: z.string().optional(),
  type: z.string(),
  phase: z.enum(['pre', 'download', 'after']),
  name: z.string(),

  error: z.string().optional(),
  inputs: z.array(NodeInput),
  status: z.enum(['pending', 'running', 'success', 'failed']).optional(),
  config: z.record(z.string(), z.unknown()).optional(),
})

export type NodeInput = z.infer<typeof NodeInput>

export type PipelineNode = z.infer<typeof PipelineNode>
