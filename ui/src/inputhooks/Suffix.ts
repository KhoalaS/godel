import z from 'zod'

const SuffixFuncArg = z.object({
  input: z.string().optional().default(''),
  suffix: z.string().optional().default(''),
})

export function SuffixHook(arg: unknown): string {
  const suffixArg = z.parse(SuffixFuncArg, arg)

  return suffixArg.input + suffixArg.suffix
}
