<script setup>
import {computed} from 'vue'
import {useRoute} from 'vue-router'
import PrototypeAvatar from './components/PrototypeAvatar.vue'
import {currentUser} from './data/mock-data'
import './styles/prototype.css'

const route = useRoute()
const isGuestScreen = computed(() => route.meta.prototypeGuest === true)

const navItems = [
	{to: '/prototype/conversations', label: 'Chat', icon: 'message'},
	{to: '/prototype/users', label: 'Persone', icon: 'users'},
	{to: '/prototype/settings', label: 'Profilo', icon: 'settings'}
]
</script>

<template>
	<div class="v2-prototype">
		<div v-if="isGuestScreen" class="v2-prototype__guest">
			<RouterView />
		</div>

		<div v-else class="v2-app-shell">
			<aside class="v2-sidebar" aria-label="Navigazione principale">
				<RouterLink to="/prototype/conversations" class="v2-brand" aria-label="WasaText, conversazioni">
					<span class="v2-brand__mark">W</span>
					<span class="v2-brand__name">WasaText</span>
				</RouterLink>

				<nav class="v2-sidebar__nav">
					<RouterLink v-for="item in navItems" :key="item.to" :to="item.to" class="v2-nav-item">
						<svg v-if="item.icon === 'message'" viewBox="0 0 24 24" aria-hidden="true"><path d="M21 15a4 4 0 0 1-4 4H8l-5 3V7a4 4 0 0 1 4-4h10a4 4 0 0 1 4 4z"/></svg>
						<svg v-else-if="item.icon === 'users'" viewBox="0 0 24 24" aria-hidden="true"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2M9 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"/></svg>
						<svg v-else viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.7 1.7 0 0 0 .34 1.88l.06.06-2.83 2.83-.06-.06a1.7 1.7 0 0 0-1.88-.34 1.7 1.7 0 0 0-1 1.55V21h-4v-.08a1.7 1.7 0 0 0-1-1.55 1.7 1.7 0 0 0-1.88.34l-.06.06-2.83-2.83.06-.06A1.7 1.7 0 0 0 4.6 15a1.7 1.7 0 0 0-1.55-1H3v-4h.08a1.7 1.7 0 0 0 1.55-1 1.7 1.7 0 0 0-.34-1.88l-.06-.06 2.83-2.83.06.06A1.7 1.7 0 0 0 9 4.6a1.7 1.7 0 0 0 1-1.55V3h4v.08a1.7 1.7 0 0 0 1 1.55 1.7 1.7 0 0 0 1.88-.34l.06-.06 2.83 2.83-.06.06A1.7 1.7 0 0 0 19.4 9a1.7 1.7 0 0 0 1.55 1H21v4h-.08a1.7 1.7 0 0 0-1.52 1z"/></svg>
						<span>{{ item.label }}</span>
					</RouterLink>
				</nav>

				<div class="v2-sidebar__profile">
					<PrototypeAvatar :initials="currentUser.initials" :src="currentUser.photoUrl" size="small" online />
					<span class="v2-sidebar__identity">
						<strong>{{ currentUser.name }}</strong>
					</span>
					<RouterLink to="/prototype/login" class="v2-logout-button" aria-label="Esci da WasaText">
						<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M10 17l5-5-5-5M15 12H3M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/></svg>
						<span>Esci</span>
					</RouterLink>
				</div>
			</aside>

			<main class="v2-main">
				<RouterView />
			</main>

			<nav class="v2-mobile-nav" aria-label="Navigazione mobile">
				<RouterLink v-for="item in navItems" :key="item.to" :to="item.to">
					<span aria-hidden="true">{{ item.icon === 'message' ? '◌' : item.icon === 'users' ? '♙' : '⚙' }}</span>
					<small>{{ item.label }}</small>
				</RouterLink>
			</nav>
		</div>
	</div>
</template>
