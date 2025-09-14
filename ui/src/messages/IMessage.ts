import z from 'zod'

export const MessageLevel = z.enum(['error', 'warn', 'info'])
export type MessageLevel = z.infer<typeof MessageLevel>

export interface IMessage<D = unknown> {
  type: string
  data: D
  level: MessageLevel
}
