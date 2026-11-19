import { afterEach, describe, expect, it, vi } from 'vitest'
import axios from './axios'
import { fetchImage, imageLimits, prepareImage, validateImageFile } from './images'
const png = new Uint8Array([137, 80, 78, 71, 13, 10, 26, 10])
const file = (changes = {}) => ({
	type: 'image/png',
	name: 'foto.png',
	size: png.length,
	arrayBuffer: async () => png.buffer,
	...changes
})
afterEach(() => {
	vi.restoreAllMocks()
	vi.unstubAllGlobals()
})

describe('image attachments', () => {
	it('rejects unsupported formats, very large originals and corrupt content before upload', async () => {
		expect(() => validateImageFile(file({ type: 'image/svg+xml' }))).toThrow('PNG o JPEG')
		expect(() => validateImageFile(file({ size: imageLimits.sourceBytes + 1 }))).toThrow('25 MB')
		await expect(prepareImage(file({ type: 'image/jpeg' }))).rejects.toThrow('danneggiato')
		await expect(
			prepareImage(
				file({
					arrayBuffer: async () => {
						throw new Error('read failure')
					}
				})
			)
		).rejects.toThrow('Impossibile leggere')
	})
	it('reduces a phone image to the API upload limits', async () => {
		const compressed = new Blob([png], { type: 'image/png' })
		const encode = vi.fn().mockResolvedValue(compressed)
		const upload = await prepareImage(
			file({ size: imageLimits.bytes + 1 }),
			async () => ({ width: 4000, height: 3000 }),
			encode
		)
		expect(encode).toHaveBeenCalledOnce()
		expect(upload).toEqual({
			mimeType: 'image/png',
			name: 'foto.png',
			data: btoa(String.fromCharCode(...png))
		})
	})
	it('checks dimensions and creates an explicit MIME/base64 upload object', async () => {
		await expect(
			prepareImage(file(), async () => ({ width: imageLimits.sourceDimension + 1, height: 1 }))
		).rejects.toThrow('troppo grande')
		const upload = await prepareImage(file(), async () => ({ width: 2, height: 3 }))
		expect(upload).toEqual({
			mimeType: 'image/png',
			name: 'foto.png',
			data: btoa(String.fromCharCode(...png))
		})
	})
	it('deduplicates concurrent authenticated media reads and rejects external URLs', async () => {
		vi.stubGlobal('localStorage', { getItem: () => null })
		const blob = new Blob(['image'])
		const get = vi.spyOn(axios, 'get').mockResolvedValue({ data: blob })
		const first = fetchImage('/media/7'),
			second = fetchImage('/media/7')
		expect(first).toBe(second)
		expect(await first).toBe(blob)
		expect(get).toHaveBeenCalledOnce()
		await expect(fetchImage('https://example.com/image.png')).rejects.toThrow('non valido')
	})
})
