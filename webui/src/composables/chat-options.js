import { updateConversationPreview } from '../services/server-state'
import { createPoller } from '../services/polling'
import { chatActions } from './chat-actions'
import { scopedClient } from '../services/request-scope'
import { apiErrorMessage } from '../services/errors'
import { prepareImage } from '../services/images'
import { createMessageTimeline } from '../services/message-timeline'

export default {
	data: function () {
		return {
			errormsg: null,
			pollError: null,
			successmsg: null,
			groupError: null,

			conversationId: this.$route.params.conversationid,
			username: localStorage.getItem('username'),
			userId: localStorage.getItem('id'),
			selectedMessageId: null,
			specificComment: null,

			loading: true,
			conversation: null,
			timeline: null,
			loadingOlder: false,
			attachment: null,
			attachmentPreview: '',
			attachmentError: null,
			sending: false,
			newMessage: '',
			forwardConversations: [],

			membersToAdd: [],

			showChatModal: false, //appare quando si clicca su Info
			confirmLeaveGroup: false, //appare quando si vuole lasciare il gruppo
			confirmDeleteMessage: false, //appare quando si vuole eliminare un messaggio
			confirmDeleteComment: false,
			showMessageOptions: false, //appare quando si clicca sui puntini del messaggio
			forwardModal: false, //appare quando si vuole inoltrare un messaggio
			addMemberModal: false, //appare quando si vuole aggiungere un membro
			emojiPicker: false, //appare quando si vuole commentare un messaggio
			allComments: false //appare quando si vuole eliminare un commento
		}
	},
	computed: {
		messages() {
			return this.timeline?.state.messages || []
		},
		selectedComments() {
			return (
				this.messages.find((message) => message.messageId === this.selectedMessageId)?.comments ||
				[]
			)
		}
	},
	watch: {
		'$route.params.conversationid'(value) {
			if (value) this.startChat()
		}
	},
	methods: {
		...chatActions,
		back() {
			this.$router.push('/conversations')
		},
		async startChat() {
			this.chatController?.abort()
			this.poller?.stop()
			this.clearAttachment()
			this.conversationId = this.$route.params.conversationid
			this.announcement = ''
			this.acting = false
			this.loading = true
			this.selectedMessageId = null
			this.specificComment = null
			this.membersToAdd = []
			this.forwardConversations = []
			this.successmsg = null
			this.groupError = null
			this.conversation = null
			this.newMessage = ''
			this.errormsg = null
			this.pollError = null
			this.loadingOlder = false
			this.sending = false
			for (const key of [
				'showChatModal',
				'confirmLeaveGroup',
				'confirmDeleteMessage',
				'confirmDeleteComment',
				'showMessageOptions',
				'forwardModal',
				'addMemberModal',
				'emojiPicker',
				'allComments'
			])
				this[key] = false
			const controller = new AbortController()
			this.chatController = controller
			this.chatClient = scopedClient(this.$axios, controller.signal)
			this.timeline = createMessageTimeline(this.$axios, this.conversationId, controller.signal)
			try {
				await Promise.all([this.getConversationInfo(), this.timeline.initial()])
				if (controller.signal.aborted) return
				updateConversationPreview(this.conversationId, this.messages)
				this.scrollToBottom()
			} catch (error) {
				if (error.code === 'ERR_CANCELED') return
				if (!controller.signal.aborted)
					this.errormsg =
						error.response?.status === 404 ? 'Conversazione non trovata' : apiErrorMessage(error)
			}
			if (!controller.signal.aborted) this.loading = false
			if (!controller.signal.aborted && this.conversation) {
				this.poller = createPoller(
					() => this.pollMessages(controller),
					() => (this.timeline?.state.hasMoreChanges ? 0 : 2000)
				)
				this.poller.start()
			}
		},
		async pollMessages(controller) {
			if (controller.signal.aborted) return
			const container =
				this.$refs.messageList || document.scrollingElement || document.documentElement
			const follow =
				!container || container.scrollHeight - container.scrollTop - container.clientHeight < 100
			await this.getChatMessages()
			if (follow && !controller.signal.aborted) this.scrollToBottom()
		},
		async getConversationInfo() {
			const id = this.conversationId
			try {
				const response = await this.chatClient.get(`/conversations/${id}`, {
					signal: this.chatController?.signal
				})
				if (id === this.conversationId && !this.chatController.signal.aborted)
					this.conversation = response.data.conversation
			} catch (error) {
				if (error.code === 'ERR_CANCELED') return
				if (id === this.conversationId && error.code !== 'ERR_CANCELED')
					this.errormsg =
						error.response?.status === 404 ? 'Conversazione non trovata' : apiErrorMessage(error)
			}
		},
		async getChatMessages() {
			const timeline = this.timeline
			try {
				await timeline?.sync()
				if (timeline === this.timeline) {
					this.pollError = null
					updateConversationPreview(this.conversationId, this.messages)
				}
			} catch (error) {
				if (error.code === 'ERR_CANCELED') return
				if (timeline === this.timeline && error.code !== 'ERR_CANCELED')
					this.pollError = apiErrorMessage(error)
			}
		},
		async loadOlder() {
			if (this.loadingOlder) return
			this.loadingOlder = true
			const timeline = this.timeline,
				controller = this.chatController
			const container =
					this.$refs.messageList || document.scrollingElement || document.documentElement,
				height = container.scrollHeight
			try {
				await timeline.older()
				if (timeline !== this.timeline || controller.signal.aborted) return
				await this.$nextTick()
				container.scrollTop += container.scrollHeight - height
			} catch (error) {
				if (error.code === 'ERR_CANCELED') return
				if (timeline === this.timeline && error.code !== 'ERR_CANCELED')
					this.errormsg = apiErrorMessage(error)
			} finally {
				if (timeline === this.timeline) this.loadingOlder = false
			}
		},
		clearAttachment() {
			if (this.attachmentPreview) URL.revokeObjectURL(this.attachmentPreview)
			this.attachment = null
			this.attachmentPreview = ''
			this.attachmentError = null
		},
		async sendMessage() {
			if (this.sending || (!this.attachment && !this.newMessage.trim())) return
			const controller = this.chatController
			const draft = this.newMessage,
				attachment = this.attachment
			this.sending = true
			this.errormsg = null
			try {
				const body = attachment
					? { type: 'image', image: attachment }
					: { type: 'text', text: draft }
				await this.chatClient.post(`/conversations/${this.conversationId}/messages`, body, {
					signal: this.chatController?.signal
				})
				if (attachment && attachment === this.attachment) this.clearAttachment()
				if (!attachment && this.newMessage === draft) this.newMessage = ''
				await this.getChatMessages()
				this.scrollToBottom()
			} catch (error) {
				if (error.code === 'ERR_CANCELED') return
				if (error.code !== 'ERR_CANCELED') this.errormsg = apiErrorMessage(error)
			} finally {
				if (controller === this.chatController) this.sending = false
			}
		},
		scrollToBottom() {
			this.$nextTick(() => {
				const list = this.$refs.messageList
				if (list) {
					list.scrollTo({ top: list.scrollHeight, behavior: 'smooth' })
				}
			})
		},
		chatSettings() {
			this.showChatModal = true
		},
		async onFileSelected(event) {
			const file = event.target.files?.[0]
			if (!file) return
			this.attachmentError = null
			const controller = this.chatController
			try {
				const attachment = await prepareImage(file)
				if (controller?.signal.aborted) return
				this.clearAttachment()
				this.attachment = attachment
				this.attachmentPreview = URL.createObjectURL(file)
				this.attachmentError = null
			} catch (error) {
				if (error.code === 'ERR_CANCELED') return
				if (!controller?.signal.aborted) this.attachmentError = error.message
			} finally {
				event.target.value = ''
			}
		}
	},
	mounted() {
		this.startChat()
	},
	beforeUnmount() {
		this.chatController?.abort()
		this.poller?.stop()
		this.clearAttachment()
	}
}
