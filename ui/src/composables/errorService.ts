import { DefaultErrorService } from '@/error/DefaultErrorService'
import { ErrorServiceKey } from '@/InjectionKeys'
import { inject } from 'vue'

export function useErrorService() {
  const errorService = inject(ErrorServiceKey) ?? DefaultErrorService
  return errorService
}
