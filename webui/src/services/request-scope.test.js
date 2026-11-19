import { describe, expect, it, vi } from 'vitest'
import { scopedClient } from './request-scope'
const deferred = () => {
	let resolve
	const promise = new Promise((r) => (resolve = r))
	return { promise, resolve }
}

describe('screen request lifetime', () => {
	it('discards an old response even if its adapter ignores cancellation', async () => {
		const waiting = deferred(),
			controller = new AbortController()
		const client = scopedClient({ get: () => waiting.promise }, controller.signal)
		const result = client.get('/users')
		controller.abort()
		waiting.resolve({ data: ['old'] })
		await expect(result).rejects.toMatchObject({ code: 'ERR_CANCELED' })
	})
	it('keeps the newest read when responses arrive in reverse order', async () => {
		const old = deferred(),
			latest = deferred()
		const client = scopedClient(
			{ get: vi.fn().mockReturnValueOnce(old.promise).mockReturnValueOnce(latest.promise) },
			new AbortController().signal
		)
		const first = client.get('/users'),
			second = client.get('/users')
		latest.resolve('latest')
		expect(await second).toBe('latest')
		old.resolve('old')
		await expect(first).rejects.toMatchObject({ code: 'ERR_CANCELED' })
	})
	it('serializes writes and cancels queued operations on navigation', async () => {
		const waiting = deferred(),
			controller = new AbortController(),
			put = vi.fn().mockReturnValue(waiting.promise)
		const client = scopedClient({ put }, controller.signal)
		const first = client.put('/name', { name: 'first' }),
			second = client.put('/name', { name: 'second' })
		await Promise.resolve()
		expect(put).toHaveBeenCalledTimes(1)
		controller.abort()
		waiting.resolve({})
		await expect(first).rejects.toMatchObject({ code: 'ERR_CANCELED' })
		await expect(second).rejects.toMatchObject({ code: 'ERR_CANCELED' })
		expect(put).toHaveBeenCalledTimes(1)
	})
})
