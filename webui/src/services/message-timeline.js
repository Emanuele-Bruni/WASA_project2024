import { reactive } from 'vue'

function timestamp(value) {
	return value.includes('.')
		? value.replace(/\.(\d+)Z$/, (_, fraction) => `.${fraction.padEnd(9, '0')}Z`)
		: value.replace(/Z$/, '.000000000Z')
}
export function compareMessages(first, second) {
	return (
		timestamp(first.createdAt).localeCompare(timestamp(second.createdAt)) ||
		first.messageId - second.messageId
	)
}
export function mergeMessages(current, updates, deletedIds = [], floor = null) {
	const removed = new Set(deletedIds)
	const messages = new Map(
		current
			.filter((message) => !removed.has(message.messageId))
			.map((message) => [message.messageId, message])
	)
	for (const message of updates) {
		if (
			!removed.has(message.messageId) &&
			(messages.has(message.messageId) || !floor || compareMessages(message, floor) >= 0)
		)
			messages.set(message.messageId, message)
	}
	return [...messages.values()].sort(compareMessages)
}

export function createMessageTimeline(client, conversationID, signal) {
	const state = reactive({
		messages: [],
		before: null,
		hasMore: false,
		cursor: null,
		hasMoreChanges: false
	})
	let queue = Promise.resolve(),
		floor = null
	const endpoint = `/conversations/${conversationID}/messages`
	function serial(operation) {
		const result = queue.then(operation)
		queue = result.catch(() => {})
		return result
	}
	async function initial() {
		const { data } = await client.get(endpoint, { params: { limit: 50 }, signal })
		if (signal?.aborted) return
		state.messages = data.messages
		state.before = data.before
		state.hasMore = data.hasMore
		state.cursor = data.cursor
		floor = state.messages[0] || null
	}
	return {
		state,
		initial: () => serial(initial),
		sync: () =>
			serial(async () => {
				if (state.cursor === null) return initial()
				for (let batch = 0; batch < 5; batch++) {
					const { data } = await client.get(endpoint, {
						params: { since: state.cursor, limit: 100 },
						signal
					})
					if (signal?.aborted) return
					state.messages = mergeMessages(
						state.messages,
						data.messages,
						data.deletedIds,
						state.hasMore ? floor : null
					)
					state.cursor = data.cursor
					state.hasMoreChanges = data.hasMoreChanges
					if (!data.hasMoreChanges) break
				}
			}),
		older: () =>
			serial(async () => {
				if (!state.hasMore || !state.before) return
				const { data } = await client.get(endpoint, {
					params: { before: state.before, limit: 50 },
					signal
				})
				if (signal?.aborted) return
				state.messages = mergeMessages(state.messages, data.messages)
				state.before = data.before
				state.hasMore = data.hasMore
				floor = state.messages[0] || floor
			})
	}
}
