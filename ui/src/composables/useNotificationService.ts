import { NotificationServiceKey } from '@/InjectionKeys'
import { mustInject } from '@/utils/mustInject'

export function useNotificationService() {
  const notificationService = mustInject(NotificationServiceKey)

  return notificationService
}
