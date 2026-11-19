const resetHandlers = new Set()
export function onSessionClear(handler) {
	resetHandlers.add(handler)
}
export const appStorageKeys = [
	'auth',
	'authExpiresAt',
	'id',
	'username',
	'userMap',
	'lastMessages',
	'forwardedIds',
	'myPhoto'
]

let redirectToLogin = () => {}

export function setSessionRedirect(handler) {
	redirectToLogin = handler
}

export function clearAppSession(storage = localStorage) {
	appStorageKeys.forEach((key) => storage.removeItem(key))
	resetHandlers.forEach((handler) => handler())
}

export function getSessionToken(storage = localStorage) {
	const token = storage.getItem('auth')
	const expiresAt = Date.parse(storage.getItem('authExpiresAt'))
	return /^[A-Za-z0-9_-]{43}$/.test(token || '') && expiresAt > Date.now() ? token : null
}

export function saveAppSession(session, storage = localStorage) {
	if (
		!/^[A-Za-z0-9_-]{43}$/.test(session.token || '') ||
		!(Date.parse(session.expiresAt) > Date.now()) ||
		!Number.isInteger(session.user?.userId) ||
		session.user?.userId < 1 ||
		typeof session.user?.name !== 'string' ||
		!session.user.name.trim()
	) {
		throw new Error('Risposta di accesso non valida')
	}
	clearAppSession(storage)
	storage.setItem('auth', session.token)
	storage.setItem('authExpiresAt', session.expiresAt)
	storage.setItem('id', String(session.user?.userId))
	storage.setItem('username', session.user.name)
}

// A late response from an old session must not clear a more recent login.
export function endAppSession(expectedToken, storage = localStorage) {
	if (storage.getItem('auth') !== expectedToken) return
	clearAppSession(storage)
	redirectToLogin()
}
