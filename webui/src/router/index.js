import {createRouter, createWebHashHistory} from 'vue-router'
import Conversations from '../views/Conversations.vue'
import Login from '../views/Login.vue'
import Users from '../views/Users.vue'
import Settings from '../views/Settings.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/login',component: Login},
		{path: '/conversations', component: Conversations},
		{path: '/users', component: Users},
		{path: '/settings', component: Settings},
		//{path: '/some/:id/link', component: HomeView},
		//{path : '/', component:}
	]
})

export default router
