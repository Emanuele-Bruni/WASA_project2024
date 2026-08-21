<script setup>
import {computed, ref} from 'vue'
import PrototypeAvatar from '../components/PrototypeAvatar.vue'
import {conversations, currentUser, users} from '../data/mock-data'

const search = ref('')
const conversationItems = ref(conversations.map((conversation) => ({...conversation})))
const showNewGroup = ref(false)
const groupName = ref('')
const selectedUsers = ref([])
const notice = ref('')

const filteredConversations = computed(() => {
	const query = search.value.trim().toLowerCase()
	return query ? conversationItems.value.filter((item) => item.title.toLowerCase().includes(query)) : conversationItems.value
})

function createMockGroup() {
	const title = groupName.value.trim()
	const initials = title
		.split(/\s+/)
		.slice(0, 2)
		.map((word) => word.charAt(0).toUpperCase())
		.join('')
	const nextId = Math.max(...conversationItems.value.map((conversation) => conversation.id)) + 1
	const selectedMembers = users.filter((user) => selectedUsers.value.includes(user.id))

	const createdGroup = {
		id: nextId,
		title,
		initials: initials || 'NG',
		color: 'teal',
		isGroup: true,
		lastMessage: 'Il gruppo è stato creato',
		lastMessageAuthor: 'Tu',
		time: 'Ora',
		unread: 0,
		members: [currentUser, ...selectedMembers]
	}
	conversationItems.value.unshift(createdGroup)
	conversations.unshift(createdGroup)
	notice.value = `Il gruppo “${title}” è stato creato con successo.`
	showNewGroup.value = false
	groupName.value = ''
	selectedUsers.value = []
}
</script>

<template>
	<section class="v2-page v2-conversations-page">
		<header class="v2-page-header">
			<div>
				<h1>Conversazioni</h1>
			</div>
			<button class="v2-button v2-button--primary" @click="showNewGroup = true">
				<span aria-hidden="true">＋</span> Nuovo gruppo
			</button>
		</header>

		<div v-if="notice" class="v2-alert v2-alert--success" role="status">
			<span aria-hidden="true">✓</span>{{ notice }}
			<button class="v2-alert__close" aria-label="Chiudi avviso" @click="notice = ''">×</button>
		</div>

		<div class="v2-toolbar">
			<label class="v2-search">
				<span aria-hidden="true">⌕</span>
				<span class="v2-sr-only">Cerca nelle conversazioni</span>
				<input v-model="search" type="search" placeholder="Cerca una conversazione…">
			</label>
		</div>

		<div v-if="filteredConversations.length" class="v2-conversation-list">
			<RouterLink
				v-for="conversation in filteredConversations"
				:key="conversation.id"
				:to="`/prototype/conversations/${conversation.id}`"
				class="v2-conversation-card"
			>
				<PrototypeAvatar :initials="conversation.initials" :color="conversation.color" size="large" :online="conversation.online" />
				<span class="v2-conversation-card__body">
					<span class="v2-conversation-card__topline">
						<strong>{{ conversation.title }}</strong>
						<time>{{ conversation.time }}</time>
					</span>
					<span class="v2-conversation-card__message">
						<span v-if="conversation.isGroup"><strong>{{ conversation.lastMessageAuthor }}:</strong> </span>{{ conversation.lastMessage }}
					</span>
				</span>
				<span v-if="conversation.unread" class="v2-unread" :aria-label="`${conversation.unread} messaggi non letti`">{{ conversation.unread }}</span>
				<span v-else class="v2-card-arrow" aria-hidden="true">›</span>
			</RouterLink>
		</div>

		<div v-else class="v2-empty-state">
			<span class="v2-empty-state__icon" aria-hidden="true">⌕</span>
			<h2>Nessuna conversazione trovata</h2>
			<p>Prova con un nome diverso.</p>
		</div>

		<div v-if="showNewGroup" class="v2-modal-backdrop" @click.self="showNewGroup = false">
			<section class="v2-modal" role="dialog" aria-modal="true" aria-labelledby="new-group-title">
				<header class="v2-modal__header">
					<div><p class="v2-eyebrow">Nuova conversazione</p><h2 id="new-group-title">Crea un gruppo</h2></div>
					<button class="v2-icon-button" aria-label="Chiudi" @click="showNewGroup = false">×</button>
				</header>
				<form class="v2-form" @submit.prevent="createMockGroup">
					<label for="prototype-group-name">Nome del gruppo</label>
					<input id="prototype-group-name" v-model.trim="groupName" placeholder="Es. Progetto università" required maxlength="25">
					<fieldset class="v2-member-picker">
						<legend>Partecipanti <span>{{ selectedUsers.length }} selezionati</span></legend>
						<label v-for="user in users" :key="user.id" class="v2-member-option">
							<PrototypeAvatar :initials="user.initials" :color="user.color" size="small" />
							<span><strong>{{ user.name }}</strong></span>
							<input v-model="selectedUsers" type="checkbox" :value="user.id">
						</label>
					</fieldset>
					<footer class="v2-modal__actions">
						<button type="button" class="v2-button v2-button--secondary" @click="showNewGroup = false">Annulla</button>
						<button class="v2-button v2-button--primary" :disabled="selectedUsers.length < 2">Crea gruppo</button>
					</footer>
				</form>
			</section>
		</div>
	</section>
</template>
