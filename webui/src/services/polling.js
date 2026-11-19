// One pending operation per poller. stop() invalidates even clients ignoring AbortSignal.
export function createPoller(
	operation,
	delay = 5000,
	schedule = setTimeout,
	cancel = clearTimeout
) {
	let controller, timer, pending
	async function refresh() {
		if (!controller || controller.signal.aborted) return
		if (pending) return pending
		const signal = controller.signal
		pending = Promise.resolve()
			.then(() => operation(signal))
			.finally(() => {
				pending = null
			})
		return pending
	}
	async function tick() {
		const current = controller
		try {
			await refresh()
		} catch (error) {
			if (!current.signal.aborted) throw error
		} finally {
			if (controller === current && !current.signal.aborted)
				timer = schedule(tick, typeof delay === 'function' ? delay() : delay)
		}
	}
	return {
		refresh,
		start() {
			if (controller && !controller.signal.aborted) return
			controller = new AbortController()
			return tick()
		},
		stop() {
			controller?.abort()
			cancel(timer)
		}
	}
}
