import z from 'zod'
import { PipelineNode } from './Node'

export const Edge = z.object({
  id: z.string(),
  target: z.string(),
  source: z.string(),
  sourceHandle: z.string().nullable(),
  targetHandle: z.string().nullable(),
})

export const Graph = z.object({
  edges: z.array(Edge),
  nodes: z.record(z.string(), PipelineNode),
  Incoming: z.record(z.string(), z.array(PipelineNode)),
  Outgoing: z.record(z.string(), z.array(PipelineNode)),
})

export type Edge = z.infer<typeof Edge>
