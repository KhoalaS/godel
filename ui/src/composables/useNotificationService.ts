import type { NotificationService } from '@/services/NotificationService'
import { NotificationServiceImpl } from '@/services/NotificationServiceImpl'
import { Container } from '@n8n/di'

export function useNotificationService(): NotificationService {
  const notificationService = Container.get(NotificationServiceImpl)

  return notificationService
}
