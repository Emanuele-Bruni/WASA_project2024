import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Login from '../views/Login.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/home', component: HomeView},
		//{path: '/link1', component: HomeView},
		//{path: '/link2', component: HomeView},
		//{path: '/some/:id/link', component: HomeView},
		{path: '/login',component: Login}
		//{path : '/', component:}
	]
})

export default router
