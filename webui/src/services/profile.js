import { reactive } from 'vue'
import { getSessionToken, onSessionClear } from './session'
let revision = 0
export const profile = reactive({ name: '', photo: '', userId: null })
export function initials(name = '') {
	return name
		.trim()
		.split(/\s+/)
		.slice(0, 2)
		.map((part) => part[0] || '')
		.join('')
		.toUpperCase()
}
export function setProfile(user) {
	revision++
	Object.assign(profile, user)
}
onSessionClear(() => {
	revision++
	Object.assign(profile, { name: '', photo: '', userId: null })
})
export async function loadProfile(client) {
	const token = getSessionToken(),
		id = localStorage.getItem('id'),
		current = revision
	if (!token) return
	profile.name = localStorage.getItem('username') || ''
	profile.userId = Number(id)
	const { data } = await client.get(`/settings/photo/${id}`)
	if (token === getSessionToken() && current === revision) setProfile(data.user)
}
