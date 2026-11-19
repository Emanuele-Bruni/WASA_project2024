import { describe, expect, it, vi } from 'vitest'
import { compareMessages, createMessageTimeline, mergeMessages } from './message-timeline'
const message = (id, text = String(id)) => ({
	messageId: id,
	text,
	createdAt: `2026-11-08T12:00:${String(id).padStart(2, '0')}Z`
})
const response = (messages, extra = {}) => ({
	data: {
		messages,
		deletedIds: [],
		cursor: '20',
		before: 'older',
		hasMore: true,
		hasMoreChanges: false,
		...extra
	}
})

describe('incremental message timeline', () => {
	it('orders timestamps precisely, including nanoseconds and second boundaries', () => {
		expect(
			compareMessages(
				{ messageId: 1, createdAt: '2026-11-08T12:00:00Z' },
				{ messageId: 2, createdAt: '2026-11-08T12:00:00.000000001Z' }
			)
		).toBeLessThan(0)
	})
	it('merges reactions and deletions without inserting an unloaded older message', () => {
		const current = [message(10), message(11)]
		const result = mergeMessages(
			current,
			[message(1), { ...message(10), comments: ['👍'] }, message(12)],
			[11],
			message(10)
		)
		expect(result.map((entry) => entry.messageId)).toEqual([10, 12])
		expect(result[0].comments).toEqual(['👍'])
	})
	it('uses a change cursor for polling and keeps it when loading older pages', async () => {
		const client = {
			get: vi
				.fn()
				.mockResolvedValueOnce(response([message(10), message(11)]))
				.mockResolvedValueOnce(response([message(12)], { cursor: '21' }))
				.mockResolvedValueOnce(
					response([message(9)], { cursor: '99', before: null, hasMore: false })
				)
		}
		const timeline = createMessageTimeline(client, 1)
		await timeline.initial()
		await timeline.sync()
		await timeline.older()
		expect(client.get.mock.calls[1][1].params).toEqual({ since: '20', limit: 100 })
		expect(client.get.mock.calls[2][1].params).toEqual({ before: 'older', limit: 50 })
		expect(timeline.state.cursor).toBe('21')
		expect(timeline.state.messages.map((entry) => entry.messageId)).toEqual([9, 10, 11, 12])
	})
	it('serializes a pending poll and pagination instead of applying responses out of order', async () => {
		let resolvePoll
		const client = {
			get: vi
				.fn()
				.mockResolvedValueOnce(response([message(10)]))
				.mockImplementationOnce(
					() =>
						new Promise((resolve) => {
							resolvePoll = resolve
						})
				)
				.mockResolvedValueOnce(response([message(9)]))
		}
		const timeline = createMessageTimeline(client, 1)
		await timeline.initial()
		const poll = timeline.sync(),
			older = timeline.older()
		await vi.waitFor(() => expect(client.get).toHaveBeenCalledTimes(2))
		resolvePoll(response([message(11)], { cursor: '21' }))
		await Promise.all([poll, older])
		expect(client.get).toHaveBeenCalledTimes(3)
		expect(timeline.state.messages.map((entry) => entry.messageId)).toEqual([9, 10, 11])
	})
	it('drains additional change batches and ignores responses after a chat has been aborted', async () => {
		const controller = new AbortController()
		const client = {
			get: vi
				.fn()
				.mockResolvedValueOnce(response([message(10)]))
				.mockResolvedValueOnce(response([message(11)], { cursor: '21', hasMoreChanges: true }))
				.mockResolvedValueOnce(response([], { cursor: '22', deletedIds: [10] }))
		}
		const timeline = createMessageTimeline(client, 1, controller.signal)
		await timeline.initial()
		await timeline.sync()
		expect(timeline.state.messages.map((entry) => entry.messageId)).toEqual([11])
		expect(timeline.state.cursor).toBe('22')
		controller.abort()
		client.get.mockResolvedValueOnce(response([message(12)]))
		await timeline.sync()
		expect(timeline.state.messages.map((entry) => entry.messageId)).toEqual([11])
	})
})
