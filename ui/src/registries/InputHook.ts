import type { HookFunction } from '@/inputhooks/HookFunction'

/**
 * Registry to store hook functions.
 */
export const FunctionRegistry = new Map<string, HookFunction>()
