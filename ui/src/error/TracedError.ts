export class TracedError extends Error {
  public readonly innerError: Error
  public readonly origin: string

  constructor(err: Error, origin: string) {
    super(`[${origin}]: ${err.message}`)

    // required when extending built-ins
    Object.setPrototypeOf(this, new.target.prototype)

    this.name = 'TracedError'
    this.innerError = err
    this.origin = origin
  }

  getErrorChain(): string {
    let chain = this.origin
    let inner: Error | null = this.innerError

    while (inner) {
      if (inner instanceof TracedError) {
        chain += `>${inner.origin}`
        inner = inner.innerError
      } else {
        inner = null
      }
    }

    return chain
  }
}
