import { beforeEach, expect, it, vi } from 'vitest'
import { fetchCollection, serverState } from './server-state'
import { clearAppSession } from './session'
let values
beforeEach(() => {
	values = new Map([
		['auth', 'a'.repeat(43)],
		['authExpiresAt', new Date(Date.now() + 60000).toISOString()]
	])
	vi.stubGlobal('localStorage', {
		getItem: (key) => values.get(key) ?? null,
		removeItem: (key) => values.delete(key)
	})
})
it('does not restore previous session data from a delayed response', async () => {
	let resolve
	const pending = fetchCollection({ get: () => new Promise((r) => (resolve = r)) }, 'users')
	clearAppSession()
	resolve({ data: { users: [{ userId: 1 }] } })
	expect(await pending).toBe(null)
	expect(serverState.users).toEqual([])
})
it('clears all in-memory collections with the session', async () => {
	await fetchCollection(
		{ get: async () => ({ data: { conversations: [{ conversationId: 1 }] } }) },
		'conversations'
	)
	expect(serverState.conversations).toHaveLength(1)
	clearAppSession()
	expect(serverState.conversations).toEqual([])
})
