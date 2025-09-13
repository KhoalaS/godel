import z from 'zod'
import { FunctionRegistry } from '@/registries/InputHook'

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

FunctionRegistry.set('toBytes', (arg) => {
  const { amount, unit } = z.parse(ToBytesFuncArg, arg)
  if (amount < 0) throw new Error('amount was negative')
  return ByteUnits[unit] * amount
})
