<script setup>
import {ref} from 'vue'
import {useRouter} from 'vue-router'

const router = useRouter()
const username = ref('')
const loading = ref(false)

function enterPrototype() {
	loading.value = true
	window.setTimeout(() => router.push('/prototype/conversations'), 450)
}
</script>

<template>
	<main class="v2-login">
		<section class="v2-login__content">
			<h1>WasaText</h1>
			<p>Scegli il nome con cui vuoi partecipare alle conversazioni.</p>

			<form class="v2-login__form" @submit.prevent="enterPrototype">
				<label class="v2-sr-only" for="prototype-username">Nome utente</label>
				<input
					id="prototype-username"
					v-model.trim="username"
					autocomplete="username"
					placeholder="Nome utente"
					minlength="3"
					maxlength="16"
					required
				>
				<button class="v2-login__submit" :disabled="loading">
					<span v-if="loading" class="v2-spinner" aria-hidden="true"></span>
					{{ loading ? 'Accesso…' : 'Continua' }}
				</button>
			</form>
		</section>
	</main>
</template>
