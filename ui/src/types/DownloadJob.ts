import z from 'zod'

export const DownloadJob = z.object({
  url: z.string(),
  filename: z.string().optional(),
  id: z.string(),
  password: z.string().optional(),
  limit: z.number().nonnegative().optional(),
  configId: z.string().optional(),
  transformer: z.array(z.string()).optional(),
  bytesDownloaded: z.number().nonnegative().optional(),
  size: z.number().nonnegative().optional(),
  deleteOnCancel: z.boolean().optional(),
  status: z.literal(['idle', 'paused', 'canceled', 'downloading', 'done', 'error']).optional(),
  speed: z.float64().optional(),
  eta: z.float64().optional(),
  isParent: z.boolean().optional(),
  urls: z.array(z.string()).optional(),
  parentId: z.string().optional(),
  headers: z.record(z.string(), z.string()).optional()
})

export type DownloadJob = z.infer<typeof DownloadJob>
