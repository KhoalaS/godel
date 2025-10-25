import type { InjectionKey } from 'vue'
import type { ErrorService } from './error/ErrorService'
import type { HookFunctionService } from './inputhooks/HookFunctionService'
import type { NotificationService } from './services/NotificationService'

export const ErrorServiceKey = Symbol('ErrorService') as InjectionKey<ErrorService>
export const HookFunctionServiceKey = Symbol(
  'HookFunctionService',
) as InjectionKey<HookFunctionService>
export const NotificationServiceKey = Symbol(
  'NotificationService',
) as InjectionKey<NotificationService>
