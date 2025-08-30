export type HookFunction = (inputId: string) => void

export const FunctionRegistry = new Map<
  string,
  (...args: (string | number | boolean)[]) => string | number | boolean
>()

FunctionRegistry.set('basename', (...inputs) => {
  const input = inputs[0]
  if (typeof input != 'string') {
    return input
  }
  return input.trim().split('/').pop() ?? input
})
