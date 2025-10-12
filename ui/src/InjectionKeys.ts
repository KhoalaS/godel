import type { InjectionKey } from 'vue'
import type { ErrorService } from './error/ErrorService'
import type { HookFunctionService } from './inputhooks/HookFunctionService'

export const ErrorServiceKey = Symbol('ErrorService') as InjectionKey<ErrorService>
export const HookFunctionServiceKey = Symbol(
  'HookFunctionService',
) as InjectionKey<HookFunctionService>
