<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useRequestScope } from '../composables/use-request-scope'
import { saveAppSession } from '../services/session'
import { setProfile } from '../services/profile'
import { apiErrorMessage } from '../services/errors'
import AppAlert from '../components/AppAlert.vue'
const username = ref(''),
	loading = ref(false),
	error = ref(''),
	route = useRoute(),
	router = useRouter(),
	{ client, signal } = useRequestScope()
async function login() {
	if (loading.value) return
	loading.value = true
	error.value = ''
	try {
		const { data } = await client.post('/login', { name: username.value.trim() })
		saveAppSession(data)
		setProfile(data.user)
		const destination = router.resolve(
			typeof route.query.redirect === 'string' ? route.query.redirect : '/conversations'
		)
		await router.replace(destination.meta.requiresAuth ? destination.fullPath : '/conversations')
	} catch (e) {
		if (!signal.aborted) error.value = apiErrorMessage(e)
	} finally {
		if (!signal.aborted) loading.value = false
	}
}
</script>
<template>
	<main class="v2-login">
		<section class="v2-login__content">
			<h1>WasaText</h1>
			<p>Scegli il nome con cui vuoi partecipare alle conversazioni.</p>
			<AppAlert :message="error" />
			<form class="v2-login__form" @submit.prevent="login">
				<label class="v2-sr-only" for="username">Nome utente</label>
				<input
					id="username"
					v-model.trim="username"
					autocomplete="username"
					placeholder="Nome utente"
					minlength="3"
					maxlength="16"
					:disabled="loading"
					required
				/>
				<button class="v2-login__submit" :disabled="loading">
					<span v-if="loading" class="v2-spinner" aria-hidden="true"></span
					>{{ loading ? 'Accesso…' : 'Continua' }}
				</button>
			</form>
		</section>
	</main>
</template>
