<script setup>
import { onMounted, ref } from 'vue'
import { profile, initials, loadProfile, setProfile } from '../services/profile'
import { logout } from '../services/axios'
import { apiErrorMessage } from '../services/errors'
import { useRequestScope } from '../composables/use-request-scope'
import AppAvatar from '../components/AppAvatar.vue'
import AppAlert from '../components/AppAlert.vue'
import ImagePicker from '../components/ImagePicker.vue'
import AvatarCropper from '../components/AvatarCropper.vue'
const { client, signal } = useRequestScope(),
	name = ref(profile.name),
	error = ref(''),
	notice = ref(''),
	busy = ref(false),
	imageVersion = ref(0),
	cropping = ref(false)
onMounted(async () => {
	try {
		await loadProfile(client)
		if (!signal.aborted) name.value = profile.name
	} catch (e) {
		if (!signal.aborted) error.value = apiErrorMessage(e)
	}
})
async function save(url, body) {
	if (busy.value) return false
	busy.value = true
	error.value = ''
	notice.value = ''
	try {
		const { data } = await client.put(url, body)
		setProfile(data.user)
		localStorage.setItem('username', data.user.name)
		notice.value = 'Modifiche salvate.'
		return true
	} catch (e) {
		if (!signal.aborted) error.value = apiErrorMessage(e)
		return false
	} finally {
		if (!signal.aborted) busy.value = false
	}
}
async function savePhoto(image) {
	if (await save('/settings/photo', { image })) {
		imageVersion.value++
		cropping.value = false
	}
}
async function signOut() {
	try {
		await logout()
	} catch (e) {
		error.value = apiErrorMessage(e)
	}
}
</script>
<template>
	<section class="v2-page v2-settings-page">
		<header class="v2-page-header"><h1>Profilo</h1></header>
		<AppAlert :message="error" /><AppAlert
			:message="notice"
			success
			dismissible
			@dismiss="notice = ''"
		/>
		<div class="v2-settings-grid">
			<section class="v2-settings-card v2-profile-card">
				<div class="v2-profile-card__avatar">
					<AppAvatar :initials="initials(profile.name)" :src="profile.photo" size="xlarge" />
				</div>
				<h2>{{ profile.name }}</h2>
				<button
					v-if="profile.photo && !cropping"
					class="v2-button v2-button--secondary"
					@click="cropping = true"
				>
					Regola ritaglio
				</button>
				<AvatarCropper
					v-if="cropping && profile.photo"
					:src="profile.photo"
					:busy="busy"
					@save="savePhoto"
					@cancel="cropping = false"
				/>
				<ImagePicker
					:key="imageVersion"
					:busy="busy"
					label="Cambia immagine profilo"
					@save="savePhoto"
				/>
			</section>
			<div class="v2-settings-stack">
				<section class="v2-settings-card">
					<header>
						<div>
							<h2>Informazioni personali</h2>
							<p>Aggiorna il nome mostrato nelle conversazioni.</p>
						</div>
					</header>
					<form class="v2-form" @submit.prevent="save('/settings/username', { name: name.trim() })">
						<label for="profile-name">Nome utente</label
						><input
							id="profile-name"
							v-model="name"
							minlength="3"
							maxlength="16"
							:disabled="busy"
							required
						/>
						<div class="v2-form__actions">
							<button class="v2-button v2-button--primary" :disabled="busy">
								{{ busy ? 'Salvataggio…' : 'Salva modifiche' }}
							</button>
						</div>
					</form>
				</section>
				<section class="v2-settings-card v2-settings-card--danger">
					<div>
						<h2>Sessione</h2>
						<p>Esci da WasaText e torna alla schermata di accesso.</p>
					</div>
					<button class="v2-button v2-button--danger" :disabled="busy" @click="signOut">
						Esci
					</button>
				</section>
			</div>
		</div>
	</section>
</template>
