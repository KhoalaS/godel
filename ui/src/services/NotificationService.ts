export interface NotificationService {
  showNotification(config: NotificationConfig): void
}

export type NotificationConfig = {
  title: string
  message: string
  level: 'error' | 'warn' | 'info'
  durationInSeconds?: number
}
