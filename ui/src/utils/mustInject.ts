import { inject, type InjectionKey } from 'vue'

export function mustInject<T>(key: string | InjectionKey<T>, defaultValue?: T) {
  const injectedValue = inject(key, defaultValue)

  if (injectedValue == null) {
    const description = typeof key === 'string' ? key : key.description
    throw new Error(`Could not inject ${description}.`)
  }

  return injectedValue
}
