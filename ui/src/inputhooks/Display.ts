import { FunctionRegistry } from '@/registries/InputHook'
import z from 'zod'

const DisplayFuncArgs = z.object({
  input: z.unknown(),
})

export default function register() {
  FunctionRegistry.set('display', (arg) => {
    const { input } = z.parse(DisplayFuncArgs, arg)

    switch (typeof input) {
      case 'string':
      case 'number':
      case 'boolean':
        return input
      case 'bigint':
        return input.toString()
      case 'object':
        return String(input)
      default:
        return ''
    }
  })
}
