import type { ErrorService } from './ErrorService'

/**
 * Default service for error handling. Just logs the error to the console.
 */
export const DefaultErrorService: ErrorService = {
  handleError: (error, info) => console.error(error, info),
}
