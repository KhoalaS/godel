import { HookFunctionServiceKey } from '@/InjectionKeys'
import { inject } from 'vue'

export function useHookFunctionService() {
  const hookFunctionService = inject(HookFunctionServiceKey)
  if (hookFunctionService === undefined) {
    throw new Error('HookFunctionService was not provided')
  }

  return hookFunctionService
}
