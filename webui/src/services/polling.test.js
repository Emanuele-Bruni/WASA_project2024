import { describe, expect, it, vi } from 'vitest'
import { createPoller } from './polling'
describe('controlled polling', () => {
	it('shares in-flight refreshes and schedules only after completion', async () => {
		let resolve
		const operation = vi.fn(() => new Promise((r) => (resolve = r))),
			schedule = vi.fn()
		const poller = createPoller(operation, 5000, schedule, vi.fn())
		const started = poller.start()
		poller.refresh()
		poller.refresh()
		await Promise.resolve()
		expect(operation).toHaveBeenCalledTimes(1)
		expect(schedule).not.toHaveBeenCalled()
		resolve()
		await started
		expect(schedule).toHaveBeenCalledTimes(1)
		poller.stop()
	})
	it('aborts pending work and prevents a new timer after unmount', async () => {
		let resolve
		const schedule = vi.fn(),
			operation = vi.fn(() => new Promise((r) => (resolve = r)))
		const poller = createPoller(operation, 5000, schedule, vi.fn())
		const started = poller.start()
		await Promise.resolve()
		poller.stop()
		expect(operation.mock.calls[0][0].aborted).toBe(true)
		resolve()
		await started
		expect(schedule).not.toHaveBeenCalled()
	})
})
