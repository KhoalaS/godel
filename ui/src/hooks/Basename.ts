import { FunctionRegistry } from '@/registries/InputHook'

const slashRegex = /\/*$/

FunctionRegistry.set('basename', (arg) => {
  const input = arg['path']
  if (input == undefined || typeof input !== 'string') {
    return ''
  }

  const _input = input.replace(slashRegex, '')

  return _input.trim().split('/').pop() ?? _input
})
