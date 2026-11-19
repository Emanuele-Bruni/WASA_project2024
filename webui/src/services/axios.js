import axios from 'axios'
import { endAppSession, getSessionToken } from './session'

const instance = axios.create({
	baseURL: globalThis.__WASA_CONFIG__?.apiURL || __API_URL__,
	timeout: 1000 * 5
})

instance.interceptors.request.use((config) => {
	config.appSessionToken = localStorage.getItem('auth')
	const token = getSessionToken()
	config.headers = config.headers || {}
	delete config.headers.Authorization
	if (token && config.url !== '/login') config.headers.Authorization = `Bearer ${token}`
	return config
})

instance.interceptors.response.use(
	(response) => response,
	(error) => {
		if (error.response?.status === 401 && error.config && error.config.url !== '/login') {
			endAppSession(error.config.appSessionToken)
		}
		return Promise.reject(error)
	}
)

export async function logout() {
	const token = localStorage.getItem('auth')
	if (getSessionToken()) {
		try {
			await instance.post('/logout')
		} catch (error) {
			if (error.response?.status !== 401) throw error
		}
	}
	endAppSession(token)
}

export default instance
