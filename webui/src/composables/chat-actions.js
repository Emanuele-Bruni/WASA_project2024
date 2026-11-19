import { apiErrorMessage } from '../services/errors'
export const chatActions = {
	async setGroupName(editableTitle) {
		try {
			await this.chatClient.put(`/conversations/${this.conversationId}/name`, {
				title: editableTitle
			})

			this.successmsg = 'Nome gruppo aggiornato con successo'
			this.showChatModal = false

			await this.getConversationInfo()
			this.scrollToBottom()
		} catch (e) {
			if (e.code === 'ERR_CANCELED') return
			if (e.response) {
				const status = e.response.status
				if (status === 400) {
					this.groupError = apiErrorMessage(e)
				} else this.groupError = apiErrorMessage(e)
			} else {
				this.groupError = apiErrorMessage(e)
			}

			this.scrollToBottom()
		}
	},
	async leaveGroup() {
		try {
			await this.chatClient.delete(`/conversations/${this.conversationId}/members`)
			this.successmsg = 'Gruppo abbandonato con successo'
			this.showChatModal = false

			//torno alla lista delle conversazioni
			this.$router.push('/conversations')
		} catch (e) {
			if (e.code === 'ERR_CANCELED') return
			this.errormsg = apiErrorMessage(e)
		}
	},
	openMessageOptions(messageId) {
		this.selectedMessageId = messageId
		// prendo i commenti e li salvo
		this.showMessageOptions = true
	},
	async deleteMessage(messageId) {
		try {
			await this.chatClient.delete(`/conversations/${this.conversationId}/messages/${messageId}`)
			this.confirmDeleteMessage = false
			this.showMessageOptions = false
			await this.getChatMessages()
		} catch (e) {
			if (e.code === 'ERR_CANCELED') return
			if (e.response) {
				const status = e.response.status
				if (status === 403) {
					this.errormsg = "Un messaggio può essere rimosso solo dall'autore"
				} else {
					this.errormsg = apiErrorMessage(e)
				}
			} else this.errormsg = apiErrorMessage(e)
		}
	},
	async openForwardModal(messageId) {
		this.selectedMessageId = messageId
		try {
			const response = await this.chatClient.get('/conversations')
			if (response) {
				this.forwardConversations = response.data.conversations || []
				this.forwardModal = true
			}
		} catch (e) {
			if (e.code === 'ERR_CANCELED') return
			this.errormsg = apiErrorMessage(e)
		}
	},
	async forwardMessage(convId) {
		try {
			await this.chatClient.post(
				`/conversations/${this.conversationId}/messages/${this.selectedMessageId}/forward/${convId}`,
				null
			)
			await this.getChatMessages()
			this.successmsg = 'Messaggio inoltrato'
			this.forwardModal = false
			this.showMessageOptions = false
			this.scrollToBottom()
		} catch (e) {
			if (e.code === 'ERR_CANCELED') return
			this.errormsg = apiErrorMessage(e)
		}
	},
	async openAddMember() {
		try {
			const response = await this.chatClient.get('/users')
			if (response) {
				// Combino Ids e Names
				this.membersToAdd = (response.data.users || []).map((user) => ({
					id: user.userId,
					name: user.name
				}))
				this.addMemberModal = true
			}
		} catch (e) {
			if (e.code === 'ERR_CANCELED') return
			this.errormsg = apiErrorMessage(e)
		}
	},
	async addToGroup(userId) {
		try {
			await this.chatClient.put(`/conversations/${this.conversationId}/members`, { userId: userId })
			this.successmsg = 'Utente aggiunto con successo'

			await this.getConversationInfo()
		} catch (e) {
			if (e.code === 'ERR_CANCELED') return
			this.errormsg = apiErrorMessage(e)
		}
	},
	async commentMessage(selectedEmoji) {
		try {
			await this.chatClient.post(
				`/conversations/${this.conversationId}/messages/${this.selectedMessageId}/comments`,
				{ emoji: selectedEmoji }
			)

			await this.getChatMessages()
			this.emojiPicker = false
			this.showMessageOptions = false
		} catch (e) {
			if (e.code === 'ERR_CANCELED') return
			this.errormsg = apiErrorMessage(e)
		}
	},
	commentOption(commentId) {
		this.confirmDeleteComment = true
		this.specificComment = commentId
	},
	async uncommentMessage(commentId) {
		try {
			await this.chatClient.delete(
				`/conversations/${this.conversationId}/messages/${this.selectedMessageId}/comments/${commentId}`
			)
			this.confirmDeleteComment = false
			this.allComments = false
			this.showMessageOptions = false
			await this.getChatMessages()
		} catch (e) {
			if (e.code === 'ERR_CANCELED') return
			if (e.response) {
				this.confirmDeleteComment = false
				const status = e.response.status
				if (status === 403) {
					this.errormsg = apiErrorMessage(e)
				} else {
					this.confirmDeleteComment = false
					this.errormsg = apiErrorMessage(e)
				}
			} else this.errormsg = apiErrorMessage(e)
		}
	},
	async setGroupPhoto(image) {
		try {
			await this.chatClient.put(`/conversations/${this.conversationId}/photo`, { image })
			this.showChatModal = false
			await this.getConversationInfo()
			this.successmsg = 'Foto aggiornata con successo!'
		} catch (e) {
			if (e.code === 'ERR_CANCELED') return
			this.errormsg = apiErrorMessage(e)
		}
	}
}
