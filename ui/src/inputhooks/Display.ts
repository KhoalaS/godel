import z from 'zod'

const DisplayFuncArgs = z.object({
  input: z.unknown(),
})

export function DisplayHook(arg: unknown): string {
  const { input } = z.parse(DisplayFuncArgs, arg)

  switch (typeof input) {
    case 'string':
    case 'number':
    case 'boolean':
      return String(input)
    case 'bigint':
      return input.toString()
    case 'object':
      return String(input)
    default:
      return ''
  }
}
