import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import axios, { logout } from './axios'
import { getSessionToken, saveAppSession, setSessionRedirect } from './session'

const token = 'a'.repeat(43)
const login = () => ({
	user: { userId: 1, name: 'alice', photo: '' },
	token,
	expiresAt: new Date(Date.now() + 60000).toISOString()
})
let redirect
beforeEach(() => {
	const values = new Map()
	vi.stubGlobal('localStorage', {
		getItem: (key) => values.get(key) ?? null,
		setItem: (key, value) => values.set(key, String(value)),
		removeItem: (key) => values.delete(key)
	})
	redirect = vi.fn()
	setSessionRedirect(redirect)
})
afterEach(() => {
	vi.unstubAllGlobals()
	setSessionRedirect(() => {})
})
const success = (config) => Promise.resolve({ status: 200, data: {}, headers: {}, config })
const rejected = (status) => (config) => Promise.reject({ response: { status }, config })

describe('HTTP session service', () => {
	it('uses the configured API URL and timeout', () => {
		expect(axios.defaults.baseURL).toBe('http://localhost:3000')
		expect(axios.defaults.timeout).toBe(5000)
	})
	it('adds the session token automatically and keeps login public', async () => {
		saveAppSession(login())
		const response = await axios.get('/users', { adapter: success })
		expect(response.config.headers.Authorization).toBe(`Bearer ${token}`)
		const publicResponse = await axios.post('/login', {}, { adapter: success })
		expect(publicResponse.config.headers.Authorization).toBeUndefined()
	})
	it('clears and redirects on 401 but preserves the session on 403', async () => {
		saveAppSession(login())
		await expect(axios.get('/users', { adapter: rejected(403) })).rejects.toBeDefined()
		expect(getSessionToken()).toBe(token)
		expect(redirect).not.toHaveBeenCalled()
		await expect(axios.get('/users', { adapter: rejected(401) })).rejects.toBeDefined()
		expect(getSessionToken()).toBeNull()
		expect(redirect).toHaveBeenCalledOnce()
	})
	it('does not clear a newer session on an old 401', async () => {
		saveAppSession(login())
		await expect(
			axios.get('/users', {
				adapter: async (config) => {
					saveAppSession({ ...login(), token: 'b'.repeat(43) })
					return rejected(401)(config)
				}
			})
		).rejects.toBeDefined()
		expect(getSessionToken()).toBe('b'.repeat(43))
		expect(redirect).not.toHaveBeenCalled()
	})
	it('revokes the server session before clearing local state on logout', async () => {
		saveAppSession(login())
		const original = axios.defaults.adapter
		const adapter = vi.fn((config) => {
			expect(getSessionToken()).toBe(token)
			return success(config)
		})
		axios.defaults.adapter = adapter
		try {
			await logout()
		} finally {
			axios.defaults.adapter = original
		}
		expect(adapter.mock.calls[0][0].url).toBe('/logout')
		expect(adapter.mock.calls[0][0].headers.Authorization).toBe(`Bearer ${token}`)
		expect(getSessionToken()).toBeNull()
		expect(redirect).toHaveBeenCalledOnce()
	})
	it('keeps the session if server logout fails, allowing a retry', async () => {
		saveAppSession(login())
		const original = axios.defaults.adapter
		axios.defaults.adapter = rejected(500)
		try {
			await expect(logout()).rejects.toBeDefined()
		} finally {
			axios.defaults.adapter = original
		}
		expect(getSessionToken()).toBe(token)
		expect(redirect).not.toHaveBeenCalled()
	})
})
