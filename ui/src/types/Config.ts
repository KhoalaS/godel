import z from 'zod'

export const Config = z.object({
  id: z.string(),
  name: z.string(),
  destPath: z.string(),
  transformer: z.array(z.string()),
  limit: z.number().gt(0).optional(),
  deleteOnCancel: z.boolean(),
})

export type Config = z.infer<typeof Config>
