import z from 'zod'

const ByteUnitEnum = z.enum(['B', 'KB', 'MB', 'GB'])
type ByteUnitEnum = z.infer<typeof ByteUnitEnum>

const ByteUnits: Record<ByteUnitEnum, number> = {
  B: 1,
  KB: 1024,
  MB: 1024 * 1024,
  GB: 1024 * 1024 * 1024,
}

const ToBytesFuncArg = z.object({
  amount: z.number(),
  unit: ByteUnitEnum,
})

export function ToBytesHook(arg: unknown): number {
  const { amount, unit } = z.parse(ToBytesFuncArg, arg)

  return ByteUnits[unit] * amount
}
