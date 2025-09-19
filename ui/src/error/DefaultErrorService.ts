import type { ErrorService } from './ErrorService'

export const DefaultErrorService: ErrorService = {
  handleError: (error, info) => console.error(error, info),
}
