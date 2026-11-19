import { reactive } from 'vue'
import { getSessionToken, onSessionClear } from './session'

export const serverState = reactive({ conversations: [], users: [] })
let owner
export function resetServerState() {
	owner = undefined
	serverState.conversations = []
	serverState.users = []
}
export async function fetchCollection(client, collection, signal) {
	const token = getSessionToken()
	if (owner !== token) {
		resetServerState()
		owner = token
	}
	const { data } = await client.get(`/${collection}`, { signal })
	if (signal?.aborted || token !== getSessionToken()) return null
	serverState[collection] = data[collection] || []
	return serverState[collection]
}

onSessionClear(resetServerState)

export function updateConversationPreview(id, messages) {
	const item = serverState.conversations.find(
		(conversation) => String(conversation.conversationId) === String(id)
	)
	if (!item) return
	const last = messages[messages.length - 1]
	item.lastMessage = last
		? { ...last, text: last.type === 'image' ? '' : last.text?.slice(0, 200) }
		: null
	serverState.conversations.sort(
		(a, b) =>
			(b.lastMessage?.createdAt || '').localeCompare(a.lastMessage?.createdAt || '') ||
			b.conversationId - a.conversationId
	)
}
