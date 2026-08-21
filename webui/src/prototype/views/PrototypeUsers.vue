<script setup>
import {computed, ref} from 'vue'
import {useRouter} from 'vue-router'
import PrototypeAvatar from '../components/PrototypeAvatar.vue'
import {conversations, currentUser, messagesByConversation, users} from '../data/mock-data'

const router = useRouter()
const search = ref('')
const filteredUsers = computed(() => {
	const query = search.value.trim().toLowerCase()
	return query ? users.filter((user) => user.name.toLowerCase().includes(query)) : users
})

function openConversation(user) {
	const existingConversation = conversations.find((conversation) => (
		!conversation.isGroup && conversation.members.some((member) => member.id === user.id)
	))

	if (existingConversation) {
		router.push(`/prototype/conversations/${existingConversation.id}`)
		return
	}

	const nextId = Math.max(...conversations.map((conversation) => conversation.id)) + 1
	const newConversation = {
		id: nextId,
		title: user.name,
		initials: user.initials,
		color: user.color,
		isGroup: false,
		lastMessage: 'Nessun messaggio',
		lastMessageAuthor: '',
		time: 'Ora',
		unread: 0,
		online: user.status === 'Online',
		members: [currentUser, user]
	}
	conversations.unshift(newConversation)
	messagesByConversation[nextId] = []
	router.push(`/prototype/conversations/${nextId}`)
}
</script>

<template>
	<section class="v2-page">
		<header class="v2-page-header">
			<div><h1>Persone</h1></div>
		</header>

		<label class="v2-search v2-search--standalone">
			<span aria-hidden="true">⌕</span><span class="v2-sr-only">Cerca persone</span>
			<input v-model="search" type="search" placeholder="Cerca per nome…">
		</label>

		<div class="v2-people-grid">
			<article v-for="user in filteredUsers" :key="user.id" class="v2-person-card">
				<PrototypeAvatar :initials="user.initials" :color="user.color" size="xlarge" :online="user.status === 'Online'" />
				<div><h2>{{ user.name }}</h2></div>
				<button class="v2-button v2-button--secondary" @click="openConversation(user)">Scrivi</button>
			</article>
		</div>

		<div v-if="!filteredUsers.length" class="v2-empty-state">
			<span class="v2-empty-state__icon" aria-hidden="true">♙</span>
			<h2>Nessuna persona trovata</h2><p>Controlla il nome e riprova.</p>
		</div>
	</section>
</template>
