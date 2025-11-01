import type { HookFunctionService } from '@/inputhooks/HookFunctionService'
import { HookFunctionServiceImpl } from '@/inputhooks/HookFunctionServiceImpl'
import { Container } from '@n8n/di'

export function useHookFunctionService(): HookFunctionService {
  const hookFunctionService = Container.get(HookFunctionServiceImpl)
  return hookFunctionService
}
