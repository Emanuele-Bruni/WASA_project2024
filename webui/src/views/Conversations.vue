<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { serverState } from '../services/server-state'
import { initials } from '../services/profile'
import { apiErrorMessage } from '../services/errors'
import { useCollection } from '../composables/use-collection'
import AppAvatar from '../components/AppAvatar.vue'
import AppAlert from '../components/AppAlert.vue'
import AppDialog from '../components/AppDialog.vue'
const { client, signal, loading, error, refresh } = useCollection('conversations'),
	router = useRouter()
const search = ref(''),
	showGroup = ref(false),
	title = ref(''),
	selected = ref([]),
	users = ref([]),
	groupError = ref(''),
	busy = ref(false),
	loadingUsers = ref(false)
const conversations = computed(() =>
	serverState.conversations.filter((item) =>
		item.title.toLowerCase().includes(search.value.trim().toLowerCase())
	)
)
async function newGroup() {
	title.value = ''
	selected.value = []
	groupError.value = ''
	showGroup.value = true
	loadingUsers.value = true
	try {
		users.value = (await client.get('/users')).data.users
	} catch (e) {
		if (!signal.aborted) groupError.value = apiErrorMessage(e)
	} finally {
		if (!signal.aborted) loadingUsers.value = false
	}
}
async function create() {
	if (busy.value) return
	busy.value = true
	groupError.value = ''
	try {
		const { data } = await client.post('/conversations', {
			title: title.value.trim(),
			userIds: selected.value
		})
		showGroup.value = false
		await router.push(`/conversations/${data.conversation.conversationId}`)
	} catch (e) {
		if (!signal.aborted) groupError.value = apiErrorMessage(e)
	} finally {
		if (!signal.aborted) busy.value = false
	}
}
</script>
<template>
	<section class="v2-page v2-conversations-page">
		<header class="v2-page-header">
			<h1>Conversazioni</h1>
			<button
				class="v2-button v2-button--primary v2-create-group-button"
				aria-label="Nuovo gruppo"
				@click="newGroup"
			>
				<span class="v2-create-group-button__icon" aria-hidden="true">＋</span>
				<span class="v2-create-group-button__label">Nuovo gruppo</span>
			</button>
		</header>
		<div class="v2-toolbar">
			<label class="v2-search"
				><span aria-hidden="true">⌕</span><span class="v2-sr-only">Cerca nelle conversazioni</span
				><input v-model="search" type="search" placeholder="Cerca una conversazione…"
			/></label>
		</div>
		<AppAlert :message="error" retry @retry="refresh" />
		<p v-if="loading && !conversations.length" role="status">Caricamento…</p>
		<div class="v2-conversation-list">
			<RouterLink
				v-for="conversation in conversations"
				:key="conversation.conversationId"
				:to="`/conversations/${conversation.conversationId}`"
				class="v2-conversation-card"
			>
				<AppAvatar
					:initials="initials(conversation.title)"
					:src="conversation.photo"
					size="large"
				/>
				<span class="v2-conversation-card__body"
					><span class="v2-conversation-card__topline"
						><strong>{{ conversation.title }}</strong
						><time>{{ conversation.lastMessage?.hours }}</time></span
					>
					<span class="v2-conversation-card__message"
						><template v-if="conversation.lastMessage"
							><strong v-if="conversation.isGroup">{{ conversation.lastMessage.author }}: </strong
							>{{
								conversation.lastMessage.type === 'image'
									? '📷 Immagine'
									: conversation.lastMessage.text
							}}</template
						><template v-else>Nessun messaggio</template></span
					></span
				><span class="v2-card-arrow" aria-hidden="true">›</span>
			</RouterLink>
		</div>
		<div v-if="!loading && !error && !conversations.length" class="v2-empty-state">
			<h2>Nessuna conversazione trovata</h2>
			<p>
				{{
					search ? 'Prova con un nome diverso.' : 'Cerca una persona o crea un gruppo per iniziare.'
				}}
			</p>
			<RouterLink to="/users" class="v2-button v2-button--secondary">Cerca persone</RouterLink>
		</div>
		<AppDialog v-if="showGroup" title="Crea un gruppo" @close="!busy && (showGroup = false)"
			><AppAlert :message="groupError" />
			<p v-if="loadingUsers" role="status">Caricamento partecipanti…</p>
			<form class="v2-form" @submit.prevent="create">
				<label for="group-name">Nome del gruppo</label
				><input
					id="group-name"
					v-model.trim="title"
					placeholder="Es. Progetto università"
					required
					maxlength="25"
					:disabled="busy"
				/>
				<fieldset class="v2-member-picker" :disabled="busy">
					<legend>
						Partecipanti <span>{{ selected.length }} selezionati</span>
					</legend>
					<label v-for="user in users" :key="user.userId" class="v2-member-option"
						><AppAvatar :initials="initials(user.name)" :src="user.photo" size="small" /><strong>{{
							user.name
						}}</strong
						><input v-model="selected" type="checkbox" :value="user.userId"
					/></label>
				</fieldset>
				<footer class="v2-modal__actions">
					<button
						type="button"
						class="v2-button v2-button--secondary"
						:disabled="busy"
						@click="showGroup = false"
					>
						Annulla</button
					><button
						class="v2-button v2-button--primary"
						:disabled="busy || loadingUsers || selected.length < 2"
					>
						{{ busy ? 'Creazione…' : 'Crea gruppo' }}
					</button>
				</footer>
			</form></AppDialog
		>
	</section>
</template>
