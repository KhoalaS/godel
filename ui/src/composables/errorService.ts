import { DefaultErrorService } from '@/error/DefaultErrorService'
import { ErrorServiceKey } from '@/InjectionKeys'
import { mustInject } from '@/utils/mustInject'

export function useErrorService() {
  const errorService = mustInject(ErrorServiceKey, DefaultErrorService)
  return errorService
}
