export type HookFunction = (inputId: string) => void

export const FunctionRegistry = new Map<
  string,
  (arg: Record<string, unknown>) => string | number | boolean
>()

const slashRegex = /\/*$/

FunctionRegistry.set('basename', (arg) => {
  const input = arg['path']
  if (input == undefined || typeof input !== 'string') {
    return ''
  }

  const _input = input.replace(slashRegex, '')

  return _input.trim().split('/').pop() ?? _input
})

const ByteUnits: Record<string, number> = {
  B: 1,
  KB: 1024,
  MB: 1024 * 1024,
  GB: 1024 * 1024 * 1024,
}

FunctionRegistry.set('toBytes', (arg) => {
  // TODO
  const amount = arg['amount']
  const unit = arg['unit']

  if (
    unit == undefined ||
    amount == undefined ||
    typeof unit !== 'string' ||
    typeof amount !== 'number'
  ) {
    return NaN
  }

  if (unit in ByteUnits) {
    return ByteUnits[unit] * amount
  }

  return NaN
})
