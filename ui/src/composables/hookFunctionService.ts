import { HookFunctionServiceKey } from '@/InjectionKeys'
import { mustInject } from '@/utils/mustInject'

export function useHookFunctionService() {
  const hookFunctionService = mustInject(HookFunctionServiceKey)

  return hookFunctionService
}
