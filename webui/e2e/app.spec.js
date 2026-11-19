import { test, expect } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'
import { randomUUID } from 'node:crypto'
const apiURL = process.env.E2E_API_URL || 'http://127.0.0.1:3000'
async function api(request, path, token, method = 'GET', data) {
	const response = await request.fetch(apiURL + path, {
		method,
		headers: token ? { Authorization: `Bearer ${token}` } : {},
		data
	})
	expect(response.ok(), `${method} ${path}: ${await response.text()}`).toBeTruthy()
	return response.status() === 204 ? {} : response.json()
}
async function seed(request) {
	const suffix = randomUUID().slice(0, 8)
	const alice = await api(request, '/login', null, 'POST', { name: 'a' + suffix })
	const bob = await api(request, '/login', null, 'POST', { name: 'b' + suffix })
	const carol = await api(request, '/login', null, 'POST', { name: 'c' + suffix })
	const first = (await api(request, `/users/${bob.user.userId}`, alice.token, 'POST', {}))
		.conversation
	const second = (await api(request, `/users/${carol.user.userId}`, alice.token, 'POST', {}))
		.conversation
	return { alice, bob, carol, first, second }
}
async function login(page, name) {
	await page.goto('/#/login')
	await page.getByLabel('Nome utente', { exact: true }).fill(name)
	await page.getByRole('button', { name: 'Continua', exact: true }).click()
	await expect(page.getByRole('heading', { name: 'Conversazioni', exact: true })).toBeVisible()
}
async function open(page, id) {
	await page.goto(`/#/conversations/${id}`)
	await expect(page.getByPlaceholder('Scrivi un messaggio…')).toBeVisible()
}
async function noOverflow(page) {
	expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true)
}

test('chat navigation, reactions, forwarding, deletion and missing IDs', async ({
	page,
	request
}) => {
	const { alice, bob, carol, first, second } = await seed(request),
		errors = []
	page.on('pageerror', (e) => errors.push(e.message))
	await login(page, alice.user.name)
	await open(page, first.conversationId)
	const field = page.getByPlaceholder('Scrivi un messaggio…')
	await field.fill('Messaggio persistente')
	await page.getByRole('button', { name: 'Invia messaggio', exact: true }).click()
	await expect(page.getByText('Messaggio persistente', { exact: true })).toBeVisible()
	await expect.poll(() => page.evaluate(() => window.scrollY)).toBe(0)
	await page.getByRole('button', { name: `Azioni sul messaggio di ${alice.user.name}` }).click()
	await page.getByRole('button', { name: 'Reagisci con 👍', exact: true }).click()
	await expect(page.getByTitle('Rimuovi la tua reazione')).toBeVisible()
	await api(request, `/conversations/${first.conversationId}/messages`, bob.token, 'POST', {
		type: 'text',
		text: 'Risposta di Bob'
	})
	await expect(page.getByRole('status').filter({ hasText: 'Nuovo messaggio da' })).toContainText(
		'Risposta di Bob'
	)

	await field.fill('Bozza precedente')
	await open(page, second.conversationId)
	await expect(field).toHaveValue('')
	await expect(page.getByText('Messaggio persistente', { exact: true })).toHaveCount(0)
	await open(page, first.conversationId)
	await expect(page.getByTitle('Rimuovi la tua reazione')).toBeVisible()
	await page.getByRole('button', { name: `Azioni sul messaggio di ${alice.user.name}` }).click()
	await page.getByRole('button', { name: '↗ Inoltra', exact: true }).click()
	await page.locator('.v2-forward-option').filter({ hasText: carol.user.name }).click()
	await expect(page.getByRole('status').filter({ hasText: 'Messaggio inoltrato' })).toBeVisible()
	const forwarded = await api(
		request,
		`/conversations/${second.conversationId}/messages`,
		carol.token
	)
	expect(forwarded.messages[0].forwarded).toBe(true)
	await page.getByRole('button', { name: `Azioni sul messaggio di ${alice.user.name}` }).click()
	await page.getByRole('button', { name: 'Elimina', exact: true }).click()
	await expect(page.getByText('Messaggio persistente', { exact: true })).toHaveCount(0)
	await page.goto('/#/conversations/999999999')
	await expect(page.getByRole('heading', { name: 'Conversazione non trovata' })).toBeVisible()
	await expect(field).toHaveCount(0)
	expect(errors).toEqual([])
})

test('long text, multiline composer and all screens have no horizontal overflow', async ({
	page,
	request
}) => {
	const { alice, first } = await seed(request)
	await api(request, `/conversations/${first.conversationId}/messages`, alice.token, 'POST', {
		type: 'text',
		text: 'https://example.test/' + 'a'.repeat(900)
	})
	await login(page, alice.user.name)
	await open(page, first.conversationId)
	await noOverflow(page)
	const field = page.getByPlaceholder('Scrivi un messaggio…'),
		height = await field.evaluate((el) => el.clientHeight)
	await field.fill(Array(15).fill('riga').join('\n'))
	await expect.poll(() => field.evaluate((el) => el.clientHeight)).toBeGreaterThan(height)
	expect(await field.evaluate((el) => el.clientHeight)).toBeLessThanOrEqual(112)
	expect(await field.evaluate((el) => el.scrollHeight > el.clientHeight)).toBe(true)
	const originalViewport = page.viewportSize()
	await page.setViewportSize({ width: 390, height: 500 })
	await field.focus()
	await expect
		.poll(async () => {
			const box = await page.locator('.v2-composer').boundingBox()
			return box.y + box.height
		})
		.toBeLessThanOrEqual(500)
	const composerBox = await page.locator('.v2-composer').boundingBox(),
		navigationBox = await page.locator('.v2-mobile-nav').boundingBox()
	expect(Math.abs(composerBox.y + composerBox.height - navigationBox.y)).toBeLessThanOrEqual(1)
	await page.setViewportSize(originalViewport)
	const touchDevice = await page.evaluate(() => matchMedia('(pointer: coarse)').matches)
	await field.fill('Invio mobile')
	await field.press('Enter')
	if (touchDevice) await expect(field).toHaveValue('Invio mobile\n')
	else {
		await expect(field).toHaveValue('')
		await expect(page.getByText('Invio mobile', { exact: true })).toBeVisible()
	}
	// Shift+Enter inserts a newline; IME Enter must not submit an unfinished composition.
	await field.fill('riga')
	await field.press('Shift+Enter')
	await expect(field).toHaveValue('riga\n')
	await field.dispatchEvent('keydown', { key: 'Enter', code: 'Enter', isComposing: true })
	await expect(field).toHaveValue('riga\n')
	for (const path of ['/conversations', '/users', '/settings']) {
		await page.goto('/#' + path)
		await page.waitForLoadState('networkidle')
		await noOverflow(page)
	}
})

test('group controls, focus trap, Escape and profile identity', async ({ page, request }) => {
	const { alice, bob, carol } = await seed(request)
	await login(page, alice.user.name)
	await page.getByRole('link', { name: 'Salta al contenuto' }).focus()
	await page.keyboard.press('Enter')
	await expect(page.locator('#main-content')).toBeFocused()
	const trigger = page.getByRole('button', { name: 'Nuovo gruppo', exact: false })
	await trigger.click()
	const dialog = page.getByRole('dialog', { name: 'Crea un gruppo' })
	await expect(dialog).toBeFocused()
	expect(await page.evaluate(() => document.body.style.overflow)).toBe('hidden')
	await page.keyboard.press('Shift+Tab')
	await expect(dialog.getByRole('button', { name: 'Annulla', exact: true })).toBeFocused()
	await page.keyboard.press('Tab')
	await expect(dialog.getByRole('button', { name: 'Chiudi', exact: true })).toBeFocused()
	await page.keyboard.press('Escape')
	await expect(dialog).toHaveCount(0)
	await expect(trigger).toBeFocused()
	await trigger.click()
	await page.getByLabel('Nome del gruppo', { exact: true }).fill('Gruppo test')
	for (const user of [bob, carol])
		await page
			.locator('.v2-member-option')
			.filter({ hasText: user.user.name })
			.locator('input')
			.check()
	await page.getByRole('button', { name: 'Crea gruppo', exact: true }).click()
	await expect(page.getByRole('heading', { name: 'Gruppo test', exact: true })).toBeVisible()
	await page.getByRole('button', { name: 'Informazioni conversazione' }).click()
	await page.getByLabel('Nome del gruppo', { exact: true }).fill('Gruppo aggiornato')
	await page.getByRole('button', { name: 'Salva nome', exact: true }).click()
	await expect(page.getByRole('heading', { name: 'Gruppo aggiornato', exact: true })).toBeVisible()
	await page.goto('/#/settings')
	await page.getByLabel('Nome utente', { exact: true }).fill('n' + randomUUID().slice(0, 8))
	await page.getByRole('button', { name: 'Salva modifiche' }).click()
	await expect(page.getByRole('status')).toContainText('Modifiche salvate')
	await page.reload()
	await expect(page.getByLabel('Nome utente', { exact: true })).toHaveValue(/^n/)
})

test('image validation, cancel, failed upload and retry preserve the draft', async ({
	page,
	request
}) => {
	const { alice, first } = await seed(request)
	await login(page, alice.user.name)
	await open(page, first.conversationId)
	const field = page.getByPlaceholder('Scrivi un messaggio…'),
		file = page.locator('.v2-composer input[type=file]')
	await field.fill('Bozza conservata')
	await file.setInputFiles({
		name: 'invalid.txt',
		mimeType: 'text/plain',
		buffer: Buffer.from('no')
	})
	await expect(page.locator('.v2-composer [role=alert]')).toContainText('PNG o JPEG')
	await expect(field).toHaveValue('Bozza conservata')
	const image = await page.evaluate(async () => {
		const canvas = document.createElement('canvas')
		canvas.width = 1800
		canvas.height = 1200
		const context = canvas.getContext('2d'),
			pixels = context.createImageData(canvas.width, canvas.height)
		for (let i = 0; i < pixels.data.length; i += 4) {
			pixels.data[i] = Math.random() * 256
			pixels.data[i + 1] = Math.random() * 256
			pixels.data[i + 2] = Math.random() * 256
			pixels.data[i + 3] = 255
		}
		context.putImageData(pixels, 0, 0)
		const blob = await new Promise((resolve) => canvas.toBlob(resolve, 'image/jpeg', 0.99)),
			bytes = new Uint8Array(await blob.arrayBuffer())
		let binary = ''
		for (let start = 0; start < bytes.length; start += 8192)
			binary += String.fromCharCode(...bytes.subarray(start, start + 8192))
		return { size: blob.size, data: btoa(binary) }
	})
	expect(image.size).toBeGreaterThan(1500000)
	const imageBuffer = Buffer.from(image.data, 'base64')
	await file.setInputFiles({ name: 'camera.jpg', mimeType: 'image/jpeg', buffer: imageBuffer })
	await page.getByRole('button', { name: 'Rimuovi allegato' }).click()
	await expect(field).toHaveValue('Bozza conservata')
	await file.setInputFiles({ name: 'camera.jpg', mimeType: 'image/jpeg', buffer: imageBuffer })
	await expect(page.getByRole('button', { name: 'Rimuovi allegato' })).toBeVisible()
	let attempt = 0
	await page.route('**/api/conversations/*/messages', (route) =>
		route.request().method() === 'POST' && attempt++ === 0
			? route.fulfill({
					status: 500,
					json: { error: { code: 'internal_error', message: 'Upload non riuscito' } }
				})
			: route.continue()
	)
	await page.getByRole('button', { name: 'Invia messaggio', exact: true }).click()
	await expect(page.getByRole('alert')).toContainText('Upload non riuscito')
	await expect(page.getByRole('button', { name: 'Rimuovi allegato' })).toBeVisible()
	await expect(field).toHaveValue('Bozza conservata')
	await page.getByRole('button', { name: 'Invia messaggio', exact: true }).click()
	await expect(page.locator('.v2-message-image')).toBeVisible()
	await expect(field).toHaveValue('Bozza conservata')
})

test('screens and open dialogs pass automated WCAG accessibility checks', async ({
	page,
	request
}) => {
	const { alice, first } = await seed(request)
	await api(request, `/conversations/${first.conversationId}/messages`, alice.token, 'POST', {
		type: 'text',
		text: 'Messaggio per il controllo di accessibilità'
	})
	async function audit() {
		const result = await new AxeBuilder({ page })
			.withTags(['wcag2a', 'wcag2aa', 'wcag21aa', 'wcag22aa'])
			.analyze()
		expect(result.violations).toEqual([])
	}
	await page.goto('/#/login')
	await audit()
	await login(page, alice.user.name)
	for (const path of [
		'/conversations',
		'/users',
		'/settings',
		`/conversations/${first.conversationId}`
	]) {
		await page.goto('/#' + path)
		await page.waitForLoadState('networkidle')
		await audit()
	}
	await page.goto('/#/conversations')
	await page.getByRole('button', { name: 'Nuovo gruppo', exact: false }).click()
	await expect(page.getByRole('dialog')).toBeFocused()
	await audit()
})

test('nested dialogs keep focus and scroll lock until the last dialog closes', async ({ page }) => {
	await page.goto('http://127.0.0.1:4173/e2e/fixtures/dialogs.html')
	const trigger = page.getByRole('button', { name: 'Apri esterno' })
	await trigger.click()
	await page.getByRole('button', { name: 'Apri interno' }).click()
	const inner = page.getByRole('dialog', { name: 'Interno' })
	await expect(inner).toBeFocused()
	await page.keyboard.press('Shift+Tab')
	await expect(page.getByRole('button', { name: 'Ultimo interno' })).toBeFocused()
	await page.keyboard.press('Tab')
	await expect(inner.getByRole('button', { name: 'Chiudi' })).toBeFocused()
	await page.keyboard.press('Escape')
	await expect(inner).toHaveCount(0)
	await expect(page.getByRole('button', { name: 'Apri interno' })).toBeFocused()
	expect(await page.evaluate(() => document.body.style.overflow)).toBe('hidden')
	await page.keyboard.press('Escape')
	await expect(trigger).toBeFocused()
	expect(await page.evaluate(() => document.body.style.overflow)).toBe('')
})
