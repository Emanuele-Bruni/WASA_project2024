<script setup>
import {computed, nextTick, ref} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import PrototypeAvatar from '../components/PrototypeAvatar.vue'
import {conversations, currentUser, getConversation, getMessages, messagesByConversation} from '../data/mock-data'

const route = useRoute()
const router = useRouter()
const conversation = computed(() => getConversation(route.params.conversationid))
const messages = ref(getMessages(route.params.conversationid).map((message) => ({
	...message,
	reactions: message.reactions.map((reaction) => ({...reaction, reactedByMe: false}))
})))
const draft = ref('')
const messageList = ref(null)
const imageInput = ref(null)
const showInfo = ref(false)
const selectedMessage = ref(null)
const messageToForward = ref(null)
const showForwardPicker = ref(false)
const showLeaveConfirmation = ref(false)
const chatNotice = ref('')

const emojis = ['👍', '💙', '😂', '✨', '😮', '🎉']

async function sendMessage() {
	const text = draft.value.trim()
	if (!text) return
	messages.value.push({
		id: Date.now(),
		authorId: currentUser.id,
		author: 'Tu',
		text,
		day: 'Oggi',
		time: new Intl.DateTimeFormat('it-IT', {hour: '2-digit', minute: '2-digit'}).format(new Date()),
		reactions: []
	})
	draft.value = ''
	await nextTick()
	messageList.value?.lastElementChild?.scrollIntoView({behavior: 'smooth'})
}

function openMessageActions(message) {
	selectedMessage.value = message
}

function addReaction(emoji) {
	if (!selectedMessage.value) return
	toggleReaction(selectedMessage.value, emoji)
	selectedMessage.value = null
}

function toggleReaction(message, emoji) {
	const existingIndex = message.reactions.findIndex((reaction) => reaction.emoji === emoji)
	if (existingIndex === -1) {
		message.reactions.push({emoji, count: 1, reactedByMe: true})
		return
	}

	const existing = message.reactions[existingIndex]
	if (existing.reactedByMe) {
		existing.count -= 1
		existing.reactedByMe = false
		if (existing.count === 0) message.reactions.splice(existingIndex, 1)
		return
	}

	existing.count += 1
	existing.reactedByMe = true
}

function openForwardPicker() {
	messageToForward.value = selectedMessage.value
	selectedMessage.value = null
	showForwardPicker.value = true
}

function forwardToConversation(targetConversation) {
	if (!messageToForward.value) return
	const targetMessages = messagesByConversation[targetConversation.id] || (messagesByConversation[targetConversation.id] = [])
	const original = messageToForward.value
	targetMessages.push({
		id: Date.now(),
		authorId: currentUser.id,
		author: 'Tu',
		text: original.text,
		imageUrl: original.imageUrl,
		imageName: original.imageName,
		forwarded: true,
		day: 'Oggi',
		time: new Intl.DateTimeFormat('it-IT', {hour: '2-digit', minute: '2-digit'}).format(new Date()),
		reactions: []
	})
	chatNotice.value = `Messaggio inoltrato a “${targetConversation.title}”.`
	showForwardPicker.value = false
	messageToForward.value = null
	window.setTimeout(() => chatNotice.value = '', 3000)
}

function leaveGroup() {
	const index = conversations.findIndex((item) => item.id === conversation.value.id)
	if (index !== -1) conversations.splice(index, 1)
	delete messagesByConversation[conversation.value.id]
	showLeaveConfirmation.value = false
	showInfo.value = false
	router.push('/prototype/conversations')
}

function chooseImage() {
	imageInput.value?.click()
}

function addImageMessage(event) {
	const file = event.target.files?.[0]
	if (!file) return

	const reader = new FileReader()
	reader.onload = async () => {
		messages.value.push({
			id: Date.now(),
			authorId: currentUser.id,
			author: 'Tu',
			text: '',
			imageUrl: reader.result,
			imageName: file.name,
			day: 'Oggi',
			time: new Intl.DateTimeFormat('it-IT', {hour: '2-digit', minute: '2-digit'}).format(new Date()),
			reactions: []
		})
		event.target.value = ''
		await nextTick()
		messageList.value?.lastElementChild?.scrollIntoView({behavior: 'smooth'})
	}
	reader.readAsDataURL(file)
}
</script>

<template>
	<section class="v2-chat-page">
		<div v-if="chatNotice" class="v2-chat-notice" role="status">✓ {{ chatNotice }}</div>
		<header class="v2-chat-header">
			<RouterLink to="/prototype/conversations" class="v2-icon-button v2-chat-header__back" aria-label="Torna alle conversazioni">‹</RouterLink>
			<PrototypeAvatar :initials="conversation.initials" :color="conversation.color" size="medium" :online="conversation.online" />
			<div class="v2-chat-header__identity">
				<h1>{{ conversation.title }}</h1>
				<p v-if="conversation.isGroup">{{ conversation.members.length }} partecipanti</p>
			</div>
			<button class="v2-icon-button" aria-label="Informazioni conversazione" @click="showInfo = true">•••</button>
		</header>

		<div ref="messageList" class="v2-message-list" aria-live="polite">
			<div v-if="!messages.length" class="v2-chat-empty">
				<span aria-hidden="true">◌</span>
				<h2>Nessun messaggio</h2>
				<p>Scrivi il primo messaggio per iniziare la conversazione.</p>
			</div>
			<template v-for="(message, index) in messages" :key="message.id">
				<div v-if="index === 0 || messages[index - 1].day !== message.day" class="v2-date-divider">
					<span>{{ message.day }}</span>
				</div>
				<article :class="['v2-message-row', {'v2-message-row--mine': message.authorId === currentUser.id}]">
				<div class="v2-message-bubble">
					<div v-if="conversation.isGroup && message.authorId !== currentUser.id" class="v2-message-bubble__author">{{ message.author }}</div>
					<div v-if="message.forwarded" class="v2-forwarded-label">↗ Inoltrato</div>
					<img v-if="message.imageUrl" class="v2-message-image" :src="message.imageUrl" :alt="message.imageName || 'Immagine inviata'">
					<p v-if="message.text">{{ message.text }}</p>
					<div class="v2-message-bubble__meta"><time>{{ message.time }}</time><span v-if="message.authorId === currentUser.id" aria-label="Consegnato">✓✓</span></div>
					<button class="v2-message-menu" :aria-label="`Azioni sul messaggio di ${message.author}`" @click="openMessageActions(message)">•••</button>
					<div v-if="message.reactions.length" class="v2-reactions">
						<button
							v-for="reaction in message.reactions"
							:key="reaction.emoji"
							:class="{'v2-reaction--mine': reaction.reactedByMe}"
							:title="reaction.reactedByMe ? 'Rimuovi la tua reazione' : 'Aggiungi questa reazione'"
							@click="toggleReaction(message, reaction.emoji)"
						>
							{{ reaction.emoji }} <span>{{ reaction.count }}</span>
						</button>
					</div>
				</div>
				</article>
			</template>
		</div>

		<form class="v2-composer" @submit.prevent="sendMessage">
			<input ref="imageInput" class="v2-sr-only" type="file" accept="image/*" tabindex="-1" @change="addImageMessage">
			<button type="button" class="v2-icon-button" aria-label="Allega immagine" @click="chooseImage">＋</button>
			<label class="v2-composer__field">
				<span class="v2-sr-only">Scrivi un messaggio</span>
				<textarea v-model="draft" rows="1" placeholder="Scrivi un messaggio…" @keydown.enter.exact.prevent="sendMessage"></textarea>
			</label>
			<button class="v2-send-button" :disabled="!draft.trim()" aria-label="Invia messaggio">
				<svg viewBox="0 0 24 24" aria-hidden="true"><path d="m22 2-7 20-4-9-9-4zM22 2 11 13"/></svg>
			</button>
		</form>

		<div v-if="selectedMessage" class="v2-action-sheet-backdrop" @click.self="selectedMessage = null">
			<section class="v2-action-sheet" role="dialog" aria-modal="true" aria-label="Azioni sul messaggio">
				<div class="v2-quick-reactions">
					<button v-for="emoji in emojis" :key="emoji" @click="addReaction(emoji)">{{ emoji }}</button>
				</div>
				<button @click="openForwardPicker"><span aria-hidden="true">↗</span> Inoltra</button>
				<button v-if="selectedMessage.authorId === currentUser.id" class="v2-danger-action" @click="messages = messages.filter((item) => item.id !== selectedMessage.id); selectedMessage = null"><span aria-hidden="true">⌫</span> Elimina</button>
				<button @click="selectedMessage = null">Annulla</button>
			</section>
		</div>

		<div v-if="showInfo" class="v2-drawer-backdrop" @click.self="showInfo = false">
			<aside class="v2-drawer" role="dialog" aria-modal="true" aria-labelledby="conversation-info-title">
				<header><h2 id="conversation-info-title">Dettagli</h2><button class="v2-icon-button" aria-label="Chiudi" @click="showInfo = false">×</button></header>
				<div class="v2-drawer__hero">
					<PrototypeAvatar :initials="conversation.initials" :color="conversation.color" size="xlarge" :online="conversation.online" />
					<h3>{{ conversation.title }}</h3><p>{{ conversation.isGroup ? 'Gruppo' : 'Conversazione privata' }}</p>
				</div>
				<section v-if="conversation.isGroup" class="v2-drawer__section">
					<h3>{{ conversation.members.length }} partecipanti</h3>
					<div v-for="member in conversation.members" :key="member.id" class="v2-member-row">
						<PrototypeAvatar :initials="member.initials" :color="member.color || 'teal'" size="small" />
						<span><strong>{{ member.name }}</strong></span>
					</div>
				</section>
				<button v-if="conversation.isGroup" class="v2-button v2-button--danger v2-button--wide" @click="showLeaveConfirmation = true">Abbandona il gruppo</button>
			</aside>
		</div>

		<div v-if="showForwardPicker" class="v2-modal-backdrop" @click.self="showForwardPicker = false">
			<section class="v2-modal v2-forward-modal" role="dialog" aria-modal="true" aria-labelledby="forward-title">
				<header class="v2-modal__header">
					<div><p class="v2-eyebrow">Inoltra messaggio</p><h2 id="forward-title">Scegli una conversazione</h2></div>
					<button class="v2-icon-button" aria-label="Chiudi" @click="showForwardPicker = false">×</button>
				</header>
				<div class="v2-forward-list">
					<button
						v-for="target in conversations.filter((item) => item.id !== conversation.id)"
						:key="target.id"
						class="v2-forward-option"
						@click="forwardToConversation(target)"
					>
						<PrototypeAvatar :initials="target.initials" :color="target.color" size="medium" />
						<span><strong>{{ target.title }}</strong><small>{{ target.isGroup ? `${target.members.length} partecipanti` : 'Conversazione privata' }}</small></span>
						<span aria-hidden="true">›</span>
					</button>
				</div>
			</section>
		</div>

		<div v-if="showLeaveConfirmation" class="v2-modal-backdrop" @click.self="showLeaveConfirmation = false">
			<section class="v2-modal v2-confirm-modal" role="alertdialog" aria-modal="true" aria-labelledby="leave-group-title">
				<div class="v2-confirm-modal__icon" aria-hidden="true">!</div>
				<h2 id="leave-group-title">Abbandonare il gruppo?</h2>
				<p>Non vedrai più “{{ conversation.title }}” nell’elenco delle conversazioni.</p>
				<footer class="v2-modal__actions">
					<button class="v2-button v2-button--secondary" @click="showLeaveConfirmation = false">Annulla</button>
					<button class="v2-button v2-button--danger" @click="leaveGroup">Abbandona</button>
				</footer>
			</section>
		</div>
	</section>
</template>
