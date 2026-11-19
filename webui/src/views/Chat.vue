<script>
import chatOptions from '../composables/chat-options'
import AppAvatar from '../components/AppAvatar.vue'
import AppAlert from '../components/AppAlert.vue'
import AppDialog from '../components/AppDialog.vue'
import ChatDetails from '../components/ChatDetails.vue'
import ChatMessage from '../components/ChatMessage.vue'
import MessageComposer from '../components/MessageComposer.vue'
import { initials } from '../services/profile'
export default {
	mixins: [chatOptions],
	components: { AppAvatar, AppAlert, AppDialog, ChatDetails, ChatMessage, MessageComposer },
	data() {
		return { acting: false, announcement: '' }
	},
	watch: {
		messages(current, previous) {
			if (!previous?.length || !current.length) return
			const latest = Math.max(...previous.map((message) => message.messageId))
			const incoming = current.filter(
				(message) => message.messageId > latest && message.authorId !== Number(this.userId)
			)
			if (incoming.length) {
				const message = incoming[incoming.length - 1]
				this.announcement = `Nuovo messaggio da ${message.author}: ${message.type === 'image' ? 'Immagine' : message.text.slice(0, 120)}`
			}
		}
	},
	computed: {
		selectedMessage() {
			return this.messages.find((message) => message.messageId === this.selectedMessageId)
		}
	},
	methods: {
		initials,
		async react(message, emoji) {
			if (this.acting) return
			this.selectedMessageId = message.messageId
			const mine = message.comments.find(
				(comment) => comment.authorId === Number(this.userId) && comment.emoji === emoji
			)
			await this.perform(
				mine ? 'uncommentMessage' : 'commentMessage',
				mine ? mine.commentId : emoji
			)
		},
		async perform(method, ...args) {
			if (this.acting) return
			const controller = this.chatController
			this.acting = true
			this.errormsg = null
			this.groupError = null
			try {
				await this[method](...args)
			} finally {
				if (controller === this.chatController) this.acting = false
			}
		},
		openForward() {
			this.showMessageOptions = false
			this.perform('openForwardModal', this.selectedMessageId)
		},
		confirmLeave() {
			this.showChatModal = false
			this.confirmLeaveGroup = true
		}
	}
}
</script>
<template>
	<section v-if="loading" class="v2-page" role="status">Caricamento…</section>
	<section v-else-if="!conversation" class="v2-page">
		<h1>{{ errormsg || 'Conversazione non trovata' }}</h1>
		<button class="v2-button v2-button--secondary" @click="startChat">Riprova</button
		><RouterLink to="/conversations" class="v2-button v2-button--primary"
			>Torna alle conversazioni</RouterLink
		>
	</section>
	<section v-else class="v2-chat-page">
		<p class="v2-sr-only" role="status" aria-live="polite" aria-atomic="true">{{ announcement }}</p>
		<header class="v2-chat-header">
			<RouterLink
				to="/conversations"
				class="v2-icon-button v2-chat-header__back"
				aria-label="Torna alle conversazioni"
				>‹</RouterLink
			>
			<AppAvatar :initials="initials(conversation.title)" :src="conversation.photo" size="medium" />
			<div class="v2-chat-header__identity">
				<h1>{{ conversation.title }}</h1>
				<p v-if="conversation.isGroup">{{ conversation.members.length }} partecipanti</p>
			</div>
			<button
				class="v2-icon-button"
				aria-label="Informazioni conversazione"
				@click="showChatModal = true"
			>
				•••
			</button>
		</header>
		<div
			ref="messageList"
			class="v2-message-list"
			role="region"
			aria-label="Cronologia messaggi"
			tabindex="0"
		>
			<AppAlert
				:message="errormsg || pollError"
				retry
				@retry="getChatMessages"
				dismissible
				@dismiss="errormsg = null"
			/><AppAlert :message="successmsg" success dismissible @dismiss="successmsg = null" />
			<button
				v-if="timeline?.state.hasMore"
				class="v2-button v2-button--secondary older-messages"
				:disabled="loadingOlder"
				@click="loadOlder"
			>
				{{ loadingOlder ? 'Caricamento…' : 'Carica messaggi precedenti' }}
			</button>
			<div v-if="!messages.length" class="v2-chat-empty">
				<h2>Nessun messaggio</h2>
				<p>Scrivi il primo messaggio per iniziare la conversazione.</p>
			</div>
			<template v-for="(message, index) in messages" :key="message.messageId"
				><div
					v-if="
						index === 0 ||
						messages[index - 1].createdAt.slice(0, 10) !== message.createdAt.slice(0, 10)
					"
					class="v2-date-divider"
				>
					<span>{{
						new Date(message.createdAt).toLocaleDateString('it-IT', {
							day: 'numeric',
							month: 'long',
							year: 'numeric'
						})
					}}</span>
				</div>
				<ChatMessage
					:message="message"
					:user-id="Number(userId)"
					:is-group="conversation.isGroup"
					@actions="openMessageOptions"
					@reaction="react"
			/></template>
		</div>
		<MessageComposer
			v-model="newMessage"
			:sending="sending"
			:attachment="attachment"
			:preview="attachmentPreview"
			:error="attachmentError"
			@send="sendMessage"
			@file="onFileSelected"
			@clear="clearAttachment"
			@dismiss-error="attachmentError = null"
		/>
		<ChatDetails
			v-if="showChatModal"
			:conversation="conversation"
			:users="addMemberModal ? membersToAdd : null"
			:busy="acting"
			:error="groupError || errormsg"
			@close="!acting && (showChatModal = false)"
			@rename="perform('setGroupName', $event)"
			@photo="perform('setGroupPhoto', $event)"
			@load-users="perform('openAddMember')"
			@add="perform('addToGroup', $event)"
			@leave="confirmLeave"
		/>
		<AppDialog
			v-else-if="confirmLeaveGroup"
			title="Abbandonare il gruppo?"
			@close="!acting && (confirmLeaveGroup = false)"
			><p>Non vedrai più “{{ conversation.title }}” nell’elenco delle conversazioni.</p>
			<AppAlert :message="errormsg" />
			<footer class="v2-modal__actions">
				<button
					class="v2-button v2-button--secondary"
					:disabled="acting"
					@click="confirmLeaveGroup = false"
				>
					Annulla</button
				><button
					class="v2-button v2-button--danger"
					:disabled="acting"
					@click="perform('leaveGroup')"
				>
					Abbandona
				</button>
			</footer></AppDialog
		>
		<AppDialog
			v-else-if="showMessageOptions && selectedMessage"
			title="Azioni sul messaggio"
			@close="!acting && (showMessageOptions = false)"
		>
			<AppAlert :message="errormsg" />
			<div class="v2-quick-reactions">
				<button
					v-for="emoji in ['👍', '💙', '😂', '✨', '😮', '🎉']"
					:key="emoji"
					:disabled="acting"
					:aria-label="`Reagisci con ${emoji}`"
					@click="perform('commentMessage', emoji)"
				>
					{{ emoji }}
				</button>
			</div>
			<div class="message-actions">
				<button class="v2-button v2-button--secondary" :disabled="acting" @click="openForward">
					↗ Inoltra</button
				><button
					v-if="selectedMessage.authorId === Number(userId)"
					class="v2-button v2-button--danger"
					:disabled="acting"
					@click="perform('deleteMessage', selectedMessageId)"
				>
					Elimina</button
				><button class="v2-button v2-button--secondary" @click="showMessageOptions = false">
					Annulla
				</button>
			</div>
			<p v-for="comment in selectedComments" :key="comment.commentId" class="reaction-detail">
				{{
					conversation.members.find((member) => member.userId === comment.authorId)?.name ||
					'Utente'
				}}: {{ comment.emoji }}
				<button
					v-if="comment.authorId === Number(userId)"
					class="v2-button v2-button--secondary"
					:disabled="acting"
					@click="perform('uncommentMessage', comment.commentId)"
				>
					Rimuovi reazione
				</button>
			</p>
		</AppDialog>
		<AppDialog
			v-else-if="forwardModal"
			title="Scegli una conversazione"
			@close="!acting && (forwardModal = false)"
			><AppAlert :message="errormsg" />
			<div class="v2-forward-list">
				<button
					v-for="target in forwardConversations.filter(
						(item) => item.conversationId !== conversation.conversationId
					)"
					:key="target.conversationId"
					class="v2-forward-option"
					:disabled="acting"
					@click="perform('forwardMessage', target.conversationId)"
				>
					<AppAvatar :initials="initials(target.title)" :src="target.photo" size="medium" /><span
						><strong>{{ target.title }}</strong
						><small>{{ target.isGroup ? 'Gruppo' : 'Conversazione privata' }}</small></span
					><span aria-hidden="true">›</span>
				</button>
			</div>
			<p v-if="forwardConversations.length < 2">
				Crea un’altra conversazione per inoltrare il messaggio.
			</p></AppDialog
		>
	</section>
</template>
