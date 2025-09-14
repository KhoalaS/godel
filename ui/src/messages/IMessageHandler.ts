import type { IMessage } from './IMessage'

export interface IMessageHandler<T extends IMessage> {
  onMessage: (msg: T) => void
}
