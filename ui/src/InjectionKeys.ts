import type { InjectionKey } from 'vue'
import type { HookFunctionService } from './inputhooks/HookFunctionService'
import type { NotificationService } from './services/NotificationService'

export const HookFunctionServiceKey = Symbol(
  'HookFunctionService',
) as InjectionKey<HookFunctionService>
export const NotificationServiceKey = Symbol(
  'NotificationService',
) as InjectionKey<NotificationService>
