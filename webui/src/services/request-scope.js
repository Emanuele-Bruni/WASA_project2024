function canceled() {
	return Object.assign(new Error('Richiesta annullata'), { code: 'ERR_CANCELED' })
}

// A scope belongs to one screen (or chat ID). Ignore responses after disposal,
// including adapters which do not implement AbortSignal. Writes stay ordered.
export function scopedClient(client, signal) {
	let writes = Promise.resolve()
	const reads = new Map()
	async function request(method, url, body, options = {}) {
		if (signal.aborted || options.signal?.aborted) throw canceled()
		const revision = (reads.get(url) || 0) + 1
		if (method === 'get') reads.set(url, revision)
		const config = { ...options, signal: options.signal || signal }
		const response = await (method === 'get' || method === 'delete'
			? client[method](url, config)
			: client[method](url, body, config))
		if (
			signal.aborted ||
			options.signal?.aborted ||
			(method === 'get' && reads.get(url) !== revision)
		)
			throw canceled()
		return response
	}
	function run(method, url, body, options) {
		if (method === 'get') return request(method, url, undefined, body)
		const result = writes.then(() =>
			request(method, url, body, method === 'delete' ? body : options)
		)
		writes = result.catch(() => {})
		return result
	}
	return Object.fromEntries(
		['get', 'delete', 'post', 'put'].map((method) => [
			method,
			(url, body, options) => run(method, url, body, options)
		])
	)
}
