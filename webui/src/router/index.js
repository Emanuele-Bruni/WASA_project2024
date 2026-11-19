import { clearAppSession, getSessionToken, setSessionRedirect } from '../services/session'
import { createMemoryHistory, createRouter, createWebHashHistory } from 'vue-router'
import Conversations from '../views/Conversations.vue'
import Login from '../views/Login.vue'
import Users from '../views/Users.vue'
import Settings from '../views/Settings.vue'
import Chat from '../views/Chat.vue'
import NotFound from '../views/NotFound.vue'

const router = createRouter({
	history:
		import.meta.env.MODE === 'test'
			? createMemoryHistory()
			: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: '/', redirect: '/login' },
		{ path: '/login', component: Login },
		{ path: '/conversations', component: Conversations, meta: { requiresAuth: true } },
		{ path: '/users', component: Users, meta: { requiresAuth: true } },
		{ path: '/settings', component: Settings, meta: { requiresAuth: true } },
		{ path: '/conversations/:conversationid', component: Chat, meta: { requiresAuth: true } },
		{ path: '/:pathMatch(.*)*', component: NotFound }
	]
})

router.beforeEach((to) => {
	if (to.meta.requiresAuth && !getSessionToken()) {
		clearAppSession()
		return { path: '/login', query: { redirect: to.fullPath } }
	}
})

setSessionRedirect(() => {
	if (router.currentRoute.value.meta.requiresAuth) return router.replace('/login')
})

export default router
