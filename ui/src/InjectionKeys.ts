import type { InjectionKey } from 'vue'
import type { ErrorService } from './error/ErrorService'

export const ErrorServiceKey = Symbol('ErrorService') as InjectionKey<ErrorService>
