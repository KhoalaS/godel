import z from 'zod'

const slashRegex = /\/*$/

const BasenameFuncArgs = z.object({
  path: z.string(),
})

export function BaseNameHook(arg: unknown): string {
  let { path } = z.parse(BasenameFuncArgs, arg)
  path = path.replace(slashRegex, '')

  return path.trim().split('/').pop() ?? path
}
