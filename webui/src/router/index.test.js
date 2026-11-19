import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { saveAppSession } from '../services/session'
import router from './index'
beforeEach(() => {
	const values = new Map()
	vi.stubGlobal('localStorage', {
		getItem: (key) => values.get(key) ?? null,
		setItem: (key, value) => values.set(key, String(value)),
		removeItem: (key) => values.delete(key)
	})
})
afterEach(() => vi.unstubAllGlobals())
describe('integrated application routing', () => {
	it('redirects the root to login', () =>
		expect(router.getRoutes().find((route) => route.path === '/').redirect).toBe('/login'))
	it('protects each application screen', () => {
		for (const path of ['/conversations', '/conversations/1', '/users', '/settings'])
			expect(router.resolve(path).meta.requiresAuth).toBe(true)
	})
	it('removes the mock prototype and resolves unknown paths to the 404 page', () => {
		expect(router.getRoutes().some((route) => route.path.startsWith('/prototype'))).toBe(false)
		expect(router.resolve('/prototype/conversations').matched[0].path).toBe('/:pathMatch(.*)*')
	})
	it('remembers the requested protected destination at login', async () => {
		await router.push('/settings')
		expect(router.currentRoute.value.path).toBe('/login')
		expect(router.currentRoute.value.query.redirect).toBe('/settings')
	})
	it('allows valid sessions', async () => {
		saveAppSession({
			user: { userId: 1, name: 'alice' },
			token: 'a'.repeat(43),
			expiresAt: new Date(Date.now() + 60000).toISOString()
		})
		await router.push('/conversations/1')
		expect(router.currentRoute.value.path).toBe('/conversations/1')
	})
})
