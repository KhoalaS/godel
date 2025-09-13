import { FunctionRegistry } from '@/registries/InputHook'
import z from 'zod'

const slashRegex = /\/*$/

const BasenameFuncArgs = z.object({
  path: z.string(),
})

export default function register() {
  FunctionRegistry.set('basename', (arg) => {
    let { path } = z.parse(BasenameFuncArgs, arg)
    path = path.replace(slashRegex, '')

    return path.trim().split('/').pop() ?? path
  })
}
