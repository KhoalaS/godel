export type HookFunction = (inputId: string) => void

export const FunctionRegistry = new Map<
  string,
  (...args: (string | number | boolean)[]) => string | number | boolean
>()

const slashRegex = /\/*$/

FunctionRegistry.set('basename', (...inputs) => {
  let input = inputs[0]
  if (typeof input != 'string') {
    return input
  }

  input = input.replace(slashRegex, '')

  return input.trim().split('/').pop() ?? input
})

const ByteUnits: Record<string, number> = {
  B: 1,
  KB: 1024,
  MB: 1024 * 1024,
  GB: 1024 * 1024 * 1024,
}

FunctionRegistry.set('toBytes', (...inputs) => {
  // TODO
  return NaN
})
