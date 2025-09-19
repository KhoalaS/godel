import type { Message } from './Message'

export interface MessageHandler<T extends Message> {
  onMessage: (msg: T) => void
}
