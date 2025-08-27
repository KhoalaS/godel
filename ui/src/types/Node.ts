import type { Edge, GraphEdge } from '@vue-flow/core'
import z from 'zod'

export const NodeIO = z.object({
  id: z.string(),
  type: z.enum(['string', 'number', 'boolean', 'directory']),
  label: z.string(),
  required: z.boolean().optional(),
  readOnly: z.boolean().optional(),
  value: z.unknown().optional(),
  options: z.array(z.string()).optional(),
})

export const PipelineNode = z.object({
  id: z.string(),
  type: z.string(),
  name: z.string(),
  nodeType: z.enum(['input', 'downloader']).optional(),
  error: z.string().optional(),
  inputs: z.record(z.string(), NodeIO).optional(),
  outputs: z.record(z.string(), NodeIO).optional(),
  status: z.enum(['pending', 'running', 'success', 'failed']).optional(),
})

export type NodeIO = z.infer<typeof NodeIO>

export type PipelineNode = z.infer<typeof PipelineNode>
