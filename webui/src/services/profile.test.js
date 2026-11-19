import { beforeEach, expect, it, vi } from 'vitest'
import { profile, initials, loadProfile, setProfile } from './profile'
import { clearAppSession } from './session'
beforeEach(() => {
	const values = new Map([
		['auth', 'a'.repeat(43)],
		['authExpiresAt', new Date(Date.now() + 60000).toISOString()],
		['id', '1'],
		['username', 'alice']
	])
	vi.stubGlobal('localStorage', {
		getItem: (key) => values.get(key) ?? null,
		removeItem: (key) => values.delete(key)
	})
	clearAppSession({ removeItem: () => {} })
})
it('derives initials from the actual profile name', () => {
	expect(initials('  Mario Rossi ')).toBe('MR')
	expect(initials('')).toBe('')
})
it('does not overwrite a saved profile with an earlier load', async () => {
	let resolve
	const pending = loadProfile({ get: () => new Promise((r) => (resolve = r)) })
	setProfile({ name: 'new', photo: '/media/2', userId: 1 })
	resolve({ data: { user: { name: 'old', photo: '/media/1', userId: 1 } } })
	await pending
	expect(profile.name).toBe('new')
	expect(profile.photo).toBe('/media/2')
})
it('clears the displayed identity at logout and discards a pending load', async () => {
	let resolve
	const pending = loadProfile({ get: () => new Promise((r) => (resolve = r)) })
	clearAppSession()
	resolve({ data: { user: { name: 'old', photo: '/media/1', userId: 1 } } })
	await pending
	expect(profile).toEqual({ name: '', photo: '', userId: null })
})
