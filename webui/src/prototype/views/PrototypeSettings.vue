<script setup>
import {ref} from 'vue'
import PrototypeAvatar from '../components/PrototypeAvatar.vue'
import {currentUser} from '../data/mock-data'

const username = ref(currentUser.name)
const photoInput = ref(null)
const saved = ref(false)

function initialsFromName(name) {
	return name
		.trim()
		.split(/\s+/)
		.filter(Boolean)
		.slice(0, 2)
		.map((part) => part.charAt(0).toUpperCase())
		.join('')
}

function saveProfile() {
	currentUser.name = username.value.trim()
	currentUser.initials = initialsFromName(currentUser.name)
	saved.value = true
	window.setTimeout(() => saved.value = false, 2500)
}

function choosePhoto() {
	photoInput.value?.click()
}

function updatePhoto(event) {
	const file = event.target.files?.[0]
	if (!file) return

	const reader = new FileReader()
	reader.onload = () => {
		currentUser.photoUrl = reader.result
		saved.value = true
		window.setTimeout(() => saved.value = false, 2500)
	}
	reader.readAsDataURL(file)
	event.target.value = ''
}

function removePhoto() {
	currentUser.photoUrl = ''
	saved.value = true
	window.setTimeout(() => saved.value = false, 2500)
}
</script>

<template>
	<section class="v2-page v2-settings-page">
		<header class="v2-page-header">
			<div><h1>Profilo</h1></div>
		</header>

		<div v-if="saved" class="v2-alert v2-alert--success" role="status"><span aria-hidden="true">✓</span>Modifiche salvate nel prototipo.</div>

		<div class="v2-settings-grid">
			<section class="v2-settings-card v2-profile-card">
				<div class="v2-profile-card__avatar">
					<PrototypeAvatar :initials="currentUser.initials" :src="currentUser.photoUrl" size="xlarge" online />
					<button class="v2-avatar-edit" aria-label="Cambia immagine profilo" @click="choosePhoto">✎</button>
				</div>
				<h2>{{ currentUser.name }}</h2>
				<input ref="photoInput" class="v2-sr-only" type="file" accept="image/*" tabindex="-1" @change="updatePhoto">
				<div class="v2-profile-card__actions">
					<button class="v2-button v2-button--secondary" @click="choosePhoto">Cambia immagine</button>
					<button class="v2-button v2-button--danger" :disabled="!currentUser.photoUrl" @click="removePhoto">Rimuovi immagine</button>
				</div>
			</section>

			<div class="v2-settings-stack">
				<section class="v2-settings-card">
					<header><div><h2>Informazioni personali</h2><p>Aggiorna il nome mostrato nelle conversazioni.</p></div></header>
					<form class="v2-form" @submit.prevent="saveProfile">
						<label for="prototype-profile-name">Nome utente</label>
						<input id="prototype-profile-name" v-model.trim="username" minlength="3" maxlength="16" required>
						<div class="v2-form__actions"><button class="v2-button v2-button--primary">Salva modifiche</button></div>
					</form>
				</section>
				<section class="v2-settings-card v2-settings-card--danger">
					<div><h2>Sessione</h2><p>Esci dal prototipo e torna alla schermata di accesso.</p></div>
					<RouterLink to="/prototype/login" class="v2-button v2-button--danger">Esci</RouterLink>
				</section>
			</div>
		</div>
	</section>
</template>
