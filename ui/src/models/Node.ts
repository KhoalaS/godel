import z from 'zod'

export const NodeIO = z.object({
  id: z.string(),
  valueType: z.literal([
    'string',
    'number',
    'boolean',
    'directory',
    'downloadjob',
    'unknown',
    'torrentservice',
    'file',
    'code',
  ]),
  label: z.string().optional(),
  required: z.boolean().optional(),
  readOnly: z.boolean().optional(),
  value: z.union([z.string(), z.number(), z.boolean()]).optional(),
  options: z.array(z.string()).optional(),
  type: z.literal(['input', 'output', 'passthrough', 'generated', 'connected_only', 'selection']),
  hooks: z.record(z.string(), z.string()).optional(),
  disabled: z.boolean().optional(),
  hookMapping: z.record(z.string(), z.string()).optional(),
})

export const PipelineNode = z.object({
  id: z.string().optional(),
  type: z.string(),
  name: z.string(),
  category: z.literal(['input', 'downloader', 'utility', 'torrent']).optional(),
  error: z.string().optional(),
  io: z.record(z.string(), NodeIO).optional(),
  status: z.literal(['pending', 'running', 'success', 'failed']).optional(),
  progress: z.float64().optional(),
})

export type NodeIO = z.infer<typeof NodeIO>

export type PipelineNode = z.infer<typeof PipelineNode>

export const HandleColors: Record<NodeIO['valueType'], string> = {
  boolean: 'green',
  directory: 'violet',
  number: 'blue',
  string: 'red',
  downloadjob: 'pink',
  unknown: 'black',
  torrentservice: 'orange',
  file: 'yellow',
  code: '#34eb8f',
}
