export function apiErrorMessage(error) {
	const message = error?.response?.data?.error?.message
	if (typeof message === 'string' && message.trim()) return message
	if (error?.code === 'ECONNABORTED') return 'Il server non ha risposto in tempo. Riprova.'
	if (!error?.response) return 'Impossibile contattare il server. Controlla la connessione.'
	return 'Impossibile completare l’operazione. Riprova.'
}
