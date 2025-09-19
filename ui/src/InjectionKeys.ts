import type { InjectionKey } from 'vue'
import type { ErrorService } from './error/ErrorService'

export const ErrorServiceKey = Symbol() as InjectionKey<ErrorService>
