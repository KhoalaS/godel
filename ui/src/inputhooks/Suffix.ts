import { FunctionRegistry } from '@/registries/InputHook'
import z from 'zod'

const SuffixFuncArg = z.object({
  input: z.string().optional().default(''),
  suffix: z.string().optional().default(''),
})

export default function register() {
  FunctionRegistry.set('suffix', (arg) => {
    const suffixArg = z.parse(SuffixFuncArg, arg)

    return suffixArg.input + suffixArg.suffix
  })
}
