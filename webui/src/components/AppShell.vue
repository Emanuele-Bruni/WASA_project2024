<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppAvatar from './AppAvatar.vue'
import AppAlert from './AppAlert.vue'
import { profile, loadProfile, initials } from '../services/profile'
import { logout } from '../services/axios'
import { useRequestScope } from '../composables/use-request-scope'
const route = useRoute(),
	error = ref(''),
	{ client } = useRequestScope()
const isGuestScreen = computed(() => !route.meta.requiresAuth)
const navItems = [
	{ to: '/conversations', label: 'Chat', icon: 'message' },
	{ to: '/users', label: 'Persone', icon: 'users' },
	{ to: '/settings', label: 'Profilo', icon: 'settings' }
]
watch(
	() => route.meta.requiresAuth,
	(value) => {
		if (value) loadProfile(client).catch(() => {})
	},
	{ immediate: true }
)
function focusContent() {
	document.getElementById('main-content')?.focus({ preventScroll: true })
}
async function signOut() {
	try {
		await logout()
	} catch {
		error.value = 'Impossibile uscire. Controlla la connessione e riprova.'
	}
}
</script>
<template>
	<div class="v2-app">
		<div v-if="isGuestScreen" class="v2-app__guest">
			<slot />
		</div>

		<div v-else class="v2-app-shell">
			<a href="#main-content" class="skip-link" @click.prevent="focusContent">Salta al contenuto</a>
			<aside class="v2-sidebar" aria-label="Navigazione principale">
				<RouterLink to="/conversations" class="v2-brand" aria-label="WasaText, conversazioni">
					<span class="v2-brand__mark">W</span>
					<span class="v2-brand__name">WasaText</span>
				</RouterLink>

				<nav class="v2-sidebar__nav">
					<RouterLink v-for="item in navItems" :key="item.to" :to="item.to" class="v2-nav-item">
						<svg v-if="item.icon === 'message'" viewBox="0 0 24 24" aria-hidden="true">
							<path d="M21 15a4 4 0 0 1-4 4H8l-5 3V7a4 4 0 0 1 4-4h10a4 4 0 0 1 4 4z" />
						</svg>
						<svg v-else-if="item.icon === 'users'" viewBox="0 0 24 24" aria-hidden="true">
							<path
								d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2M9 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"
							/>
						</svg>
						<svg v-else viewBox="0 0 24 24" aria-hidden="true">
							<circle cx="12" cy="12" r="3" />
							<path
								d="M19.4 15a1.7 1.7 0 0 0 .34 1.88l.06.06-2.83 2.83-.06-.06a1.7 1.7 0 0 0-1.88-.34 1.7 1.7 0 0 0-1 1.55V21h-4v-.08a1.7 1.7 0 0 0-1-1.55 1.7 1.7 0 0 0-1.88.34l-.06.06-2.83-2.83.06-.06A1.7 1.7 0 0 0 4.6 15a1.7 1.7 0 0 0-1.55-1H3v-4h.08a1.7 1.7 0 0 0 1.55-1 1.7 1.7 0 0 0-.34-1.88l-.06-.06 2.83-2.83.06.06A1.7 1.7 0 0 0 9 4.6a1.7 1.7 0 0 0 1-1.55V3h4v.08a1.7 1.7 0 0 0 1 1.55 1.7 1.7 0 0 0 1.88-.34l.06-.06 2.83 2.83-.06.06A1.7 1.7 0 0 0 19.4 9a1.7 1.7 0 0 0 1.55 1H21v4h-.08a1.7 1.7 0 0 0-1.52 1z"
							/>
						</svg>
						<span>{{ item.label }}</span>
					</RouterLink>
				</nav>

				<div class="v2-sidebar__profile">
					<AppAvatar :initials="initials(profile.name)" :src="profile.photo" size="small" />
					<span class="v2-sidebar__identity">
						<strong>{{ profile.name }}</strong>
					</span>
					<button
						type="button"
						class="v2-logout-button"
						aria-label="Esci da WasaText"
						@click="signOut"
					>
						<svg viewBox="0 0 24 24" aria-hidden="true">
							<path d="M10 17l5-5-5-5M15 12H3M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4" />
						</svg>
						<span>Esci</span>
					</button>
				</div>
			</aside>

			<main id="main-content" class="v2-main" tabindex="-1">
				<AppAlert v-if="error" :message="error" />
				<slot />
			</main>

			<nav class="v2-mobile-nav" aria-label="Navigazione mobile">
				<RouterLink v-for="item in navItems" :key="item.to" :to="item.to">
					<span aria-hidden="true">{{
						item.icon === 'message' ? '◌' : item.icon === 'users' ? '♙' : '⚙'
					}}</span>
					<small>{{ item.label }}</small>
				</RouterLink>
			</nav>
		</div>
	</div>
</template>
