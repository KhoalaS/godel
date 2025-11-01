import { Service } from '@n8n/di'
import type { HookFunction } from './HookFunction'
import type { HookFunctionService } from './HookFunctionService'

@Service()
export class HookFunctionServiceImpl implements HookFunctionService {
  private registry = new Map<string, HookFunction>()

  register(id: string, fn: HookFunction) {
    if (this.registry.has(id)) {
      console.warn(
        `[HookFunctionServiceImpl]: ID ${id} is already assigned. Assigning to new function.`,
      )
    }
    this.registry.set(id, fn)
  }

  executeHookFunction(id: string, arg: unknown): ReturnType<HookFunction> {
    const fn = this.registry.get(id)

    if (fn === undefined) {
      throw new Error(`No hook-function found for id ${id} in function registry.`)
    }

    return fn(arg)
  }

  getFunction(id: string): HookFunction | undefined {
    return this.registry.get(id)
  }
}
