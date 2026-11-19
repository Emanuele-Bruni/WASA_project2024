import { afterEach, describe, expect, it, vi } from 'vitest'
import {
	appStorageKeys,
	clearAppSession,
	endAppSession,
	getSessionToken,
	saveAppSession,
	setSessionRedirect
} from './session'

function memoryStorage() {
	const values = new Map()
	return {
		getItem: (key) => values.get(key) ?? null,
		setItem: (key, value) => values.set(key, String(value)),
		removeItem: (key) => values.delete(key)
	}
}
const session = () => ({
	user: { userId: 1, name: 'alice', photo: '' },
	token: 'a'.repeat(43),
	expiresAt: new Date(Date.now() + 60000).toISOString()
})
afterEach(() => setSessionRedirect(() => {}))

describe('application session', () => {
	it('clears application state without deleting unrelated preferences', () => {
		const storage = memoryStorage()
		appStorageKeys.forEach((key) => storage.setItem(key, 'value'))
		storage.setItem('unrelated', 'keep')
		clearAppSession(storage)
		appStorageKeys.forEach((key) => expect(storage.getItem(key)).toBeNull())
		expect(storage.getItem('unrelated')).toBe('keep')
	})

	it('saves a session and clears the previous user caches', () => {
		const storage = memoryStorage()
		storage.setItem('userMap', 'old user')
		const login = session()
		saveAppSession(login, storage)
		expect(getSessionToken(storage)).toBe(login.token)
		expect(storage.getItem('username')).toBe('alice')
		expect(storage.getItem('id')).toBe('1')
		expect(storage.getItem('userMap')).toBeNull()
	})

	it('rejects legacy IDs, invalid tokens and expired sessions', () => {
		const storage = memoryStorage()
		for (const token of ['1', 'Bearer 1', 'undefined', '']) {
			storage.setItem('auth', token)
			storage.setItem('authExpiresAt', session().expiresAt)
			expect(getSessionToken(storage)).toBeNull()
		}
		storage.setItem('auth', session().token)
		storage.setItem('authExpiresAt', new Date(0).toISOString())
		expect(getSessionToken(storage)).toBeNull()
		expect(() =>
			saveAppSession({ user: { userId: 1, name: 'alice', photo: '' } }, storage)
		).toThrow()
	})

	it('ignores a late unauthorized response from a previous session', () => {
		const storage = memoryStorage()
		const redirect = vi.fn()
		setSessionRedirect(redirect)
		saveAppSession(session(), storage)
		endAppSession('old token', storage)
		expect(getSessionToken(storage)).toBe(session().token)
		expect(redirect).not.toHaveBeenCalled()
		endAppSession(session().token, storage)
		expect(getSessionToken(storage)).toBeNull()
		expect(redirect).toHaveBeenCalledOnce()
	})
})
