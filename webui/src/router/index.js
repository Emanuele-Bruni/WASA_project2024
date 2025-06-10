import {createRouter, createWebHashHistory} from 'vue-router'
import Conversations from '../views/Conversations.vue'
import Login from '../views/Login.vue'
import Users from '../views/Users.vue'
import Settings from '../views/Settings.vue'
import Chat from '../views/Chat.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path : '/', redirect:'/login'},
		{path: '/login',component: Login},
		{path: '/conversations', component: Conversations},
		{path: '/users', component: Users},
		{path: '/settings', component: Settings},
		{path: '/conversations/:conversationId', component: Chat}
		//{path: '/some/:id/link', component: HomeView},
	]
})

export default router
