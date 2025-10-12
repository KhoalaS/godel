import type { HookFunction } from './HookFunction'

export interface HookFunctionService {
  register(id: string, fn: HookFunction): void
  executeHookFunction(id: string, arg: unknown): ReturnType<HookFunction>
  getFunction(id: string): HookFunction | undefined
}
