import { expect, it, vi } from 'vitest'
import chatOptions from './chat-options'
import { scopedClient } from '../services/request-scope'

it('a send completed after navigation does not clear the new chat draft', async () => {
	let resolve
	const controller = new AbortController()
	const client = scopedClient({ post: () => new Promise((r) => (resolve = r)) }, controller.signal)
	const view = {
		sending: false,
		attachment: null,
		newMessage: 'same draft',
		conversationId: '1',
		chatController: controller,
		chatClient: client,
		getChatMessages: vi.fn(),
		scrollToBottom: vi.fn()
	}
	const sending = chatOptions.methods.sendMessage.call(view)
	await Promise.resolve()
	controller.abort()
	view.chatController = new AbortController()
	view.conversationId = '2'
	view.sending = true
	resolve({ data: {} })
	await sending
	expect(view.newMessage).toBe('same draft')
	expect(view.sending).toBe(true)
	expect(view.getChatMessages).not.toHaveBeenCalled()
})

it('derives reaction details from the current timeline instead of a stale copy', () => {
	const view = { selectedMessageId: 1, messages: [{ messageId: 1, comments: [{ emoji: '👍' }] }] }
	expect(chatOptions.computed.selectedComments.call(view)).toEqual([{ emoji: '👍' }])
	view.messages = []
	expect(chatOptions.computed.selectedComments.call(view)).toEqual([])
})

it('keeps the translated missing-chat message when the timeline fails in parallel', async () => {
	const failure = { response: { status: 404, data: { error: { message: 'Not Found' } } } }
	const view = {
		$route: { params: { conversationid: '999999' } },
		$axios: { get: vi.fn().mockRejectedValue(failure) },
		clearAttachment: vi.fn(),
		scrollToBottom: vi.fn()
	}
	view.getConversationInfo = chatOptions.methods.getConversationInfo.bind(view)
	await chatOptions.methods.startChat.call(view)
	expect(view.errormsg).toBe('Conversazione non trovata')
	expect(view.loading).toBe(false)
	expect(view.conversation).toBe(null)
	view.chatController.abort()
})
