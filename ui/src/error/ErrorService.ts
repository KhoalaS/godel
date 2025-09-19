export interface ErrorService {
  handleError: ErrorHandler
}

export type ErrorHandler = (error: unknown, info?: string) => void
