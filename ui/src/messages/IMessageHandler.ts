export interface IMessageHandler<T> {
  onMessage: (msg: T) => void
}
