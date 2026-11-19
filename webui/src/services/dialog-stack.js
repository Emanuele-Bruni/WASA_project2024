// One document listener and one scroll lock, regardless of the number of dialogs.
const stack = []
let overflow, background, wasInert
function top() {
	return stack[stack.length - 1]
}
function focusables(panel) {
	return [...panel.querySelectorAll('button, input, textarea, select, a[href], [tabindex]')].filter(
		(node) =>
			!node.disabled &&
			node.tabIndex >= 0 &&
			node.getClientRects().length &&
			!node.closest('[inert]')
	)
}
function focus(event) {
	const current = top()
	if (current && !current.panel.contains(event.target)) current.panel.focus()
}
function keyboard(event) {
	const current = top()
	if (!current) return
	if (event.key === 'Escape') {
		event.preventDefault()
		event.stopImmediatePropagation()
		current.close()
		return
	}
	if (event.key !== 'Tab') return
	const nodes = focusables(current.panel),
		first = nodes[0],
		last = nodes[nodes.length - 1]
	if (!first) {
		event.preventDefault()
		current.panel.focus()
		return
	}
	if (
		event.shiftKey &&
		(document.activeElement === first || !nodes.includes(document.activeElement))
	) {
		event.preventDefault()
		last.focus()
	} else if (
		!event.shiftKey &&
		(document.activeElement === last || !current.panel.contains(document.activeElement))
	) {
		event.preventDefault()
		first.focus()
	}
}
export function registerDialog(panel, close) {
	const entry = { panel, close, previous: document.activeElement }
	if (!stack.length) {
		overflow = document.body.style.overflow
		document.body.style.overflow = 'hidden'
		background = document.querySelector('#app > .v2-app')
		wasInert = background?.inert
		if (background) background.inert = true
		document.addEventListener('keydown', keyboard, true)
		document.addEventListener('focusin', focus, true)
	} else top().panel.inert = true
	stack.push(entry)
	panel.focus()
	return () => {
		const index = stack.indexOf(entry),
			wasTop = entry === top()
		if (index < 0) return
		stack.splice(index, 1)
		const current = top()
		if (current) current.panel.inert = false
		else {
			document.body.style.overflow = overflow
			if (background) background.inert = wasInert
			document.removeEventListener('keydown', keyboard, true)
			document.removeEventListener('focusin', focus, true)
		}
		if (wasTop) {
			const previous = entry.previous
			if (
				previous?.isConnected &&
				!previous.closest('[inert]') &&
				(!current || current.panel.contains(previous))
			)
				previous.focus()
			else current?.panel.focus()
		}
	}
}
