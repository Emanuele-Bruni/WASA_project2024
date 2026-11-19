<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { serverState } from '../services/server-state'
import { initials } from '../services/profile'
import { apiErrorMessage } from '../services/errors'
import { useCollection } from '../composables/use-collection'
import AppAvatar from '../components/AppAvatar.vue'
import AppAlert from '../components/AppAlert.vue'
const { client, signal, loading, error, refresh } = useCollection('users'),
	router = useRouter(),
	search = ref(''),
	opening = ref(null)
const users = computed(() =>
	serverState.users.filter((user) =>
		user.name.toLowerCase().includes(search.value.trim().toLowerCase())
	)
)
async function open(user) {
	if (opening.value) return
	opening.value = user.userId
	try {
		const all = (await client.get('/conversations')).data.conversations
		const existing = all.find(
			(item) => !item.isGroup && item.members?.some((member) => member.userId === user.userId)
		)
		const id =
			existing?.conversationId ||
			(await client.post(`/users/${user.userId}`, {})).data.conversation.conversationId
		await router.push(`/conversations/${id}`)
	} catch (e) {
		if (!signal.aborted) error.value = apiErrorMessage(e)
	} finally {
		if (!signal.aborted) opening.value = null
	}
}
</script>
<template>
	<section class="v2-page">
		<header class="v2-page-header"><h1>Persone</h1></header>
		<label class="v2-search v2-search--standalone"
			><span aria-hidden="true">⌕</span><span class="v2-sr-only">Cerca persone</span
			><input v-model="search" type="search" placeholder="Cerca per nome…"
		/></label>
		<AppAlert :message="error" retry @retry="refresh" />
		<p v-if="loading" role="status">Caricamento…</p>
		<div class="v2-people-grid">
			<article v-for="user in users" :key="user.userId" class="v2-person-card">
				<AppAvatar :initials="initials(user.name)" :src="user.photo" size="xlarge" />
				<h2>{{ user.name }}</h2>
				<button class="v2-button v2-button--secondary" :disabled="!!opening" @click="open(user)">
					{{ opening === user.userId ? 'Apertura…' : 'Scrivi' }}
				</button>
			</article>
		</div>
		<div v-if="!loading && !error && !users.length" class="v2-empty-state">
			<h2>Nessuna persona trovata</h2>
			<p>Controlla il nome e riprova.</p>
		</div>
	</section>
</template>
