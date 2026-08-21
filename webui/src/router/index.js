import {createRouter, createWebHashHistory} from 'vue-router'
import Conversations from '../views/Conversations.vue'
import Login from '../views/Login.vue'
import Users from '../views/Users.vue'
import Settings from '../views/Settings.vue'
import Chat from '../views/Chat.vue'
import PrototypeShell from '../prototype/PrototypeShell.vue'
import PrototypeLogin from '../prototype/views/PrototypeLogin.vue'
import PrototypeConversations from '../prototype/views/PrototypeConversations.vue'
import PrototypeChat from '../prototype/views/PrototypeChat.vue'
import PrototypeUsers from '../prototype/views/PrototypeUsers.vue'
import PrototypeSettings from '../prototype/views/PrototypeSettings.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path : '/', redirect:'/login'},
		{path: '/login',component: Login},
		{path: '/conversations', component: Conversations},
		{path: '/users', component: Users},
		{path: '/settings', component: Settings},
		{path: '/conversations/:conversationid', component: Chat},
		{
			path: '/prototype',
			component: PrototypeShell,
			children: [
				{path: '', redirect: '/prototype/conversations'},
				{path: 'login', component: PrototypeLogin, meta: {prototypeGuest: true}},
				{path: 'conversations', component: PrototypeConversations},
				{path: 'conversations/:conversationid', component: PrototypeChat},
				{path: 'users', component: PrototypeUsers},
				{path: 'settings', component: PrototypeSettings}
			]
		}
		//{path: '/some/:id/link', component: HomeView},
	]
})

export default router
