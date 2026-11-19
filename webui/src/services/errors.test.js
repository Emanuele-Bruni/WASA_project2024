import { describe, expect, it } from 'vitest'
import { apiErrorMessage } from './errors'

describe('API error messages', () => {
	it('shows server validation feedback instead of a generic HTTP status', () => {
		expect(
			apiErrorMessage({
				response: {
					status: 409,
					data: { error: { code: 'conflict', message: 'Nome già utilizzato' } }
				}
			})
		).toBe('Nome già utilizzato')
	})
	it('provides useful feedback for network errors, timeouts and unexpected responses', () => {
		expect(apiErrorMessage(new Error())).toContain('connessione')
		expect(apiErrorMessage({ code: 'ECONNABORTED' })).toContain('in tempo')
		expect(apiErrorMessage({ response: { status: 500, data: '<html>' } })).toContain('Riprova')
	})
})
