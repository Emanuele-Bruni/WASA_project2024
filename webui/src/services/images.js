import axios from './axios'
import { getSessionToken } from './session'

export const imageLimits = {
	bytes: 1500000,
	dimension: 4096,
	pixels: 4194304,
	sourceBytes: 25000000,
	sourceDimension: 12000,
	sourcePixels: 50000000
}
export const imageTypes = ['image/png', 'image/jpeg']
const pending = new Map()

export function validateImageFile(file) {
	if (!file || !imageTypes.includes(file.type)) throw new Error('Scegli un’immagine PNG o JPEG.')
	if (!file.size) throw new Error('Il file immagine è vuoto.')
	if (file.size > imageLimits.sourceBytes)
		throw new Error('L’immagine originale non può superare 25 MB.')
}
async function dimensions(buffer, type) {
	const bitmap = await createImageBitmap(new Blob([buffer], { type }))
	const result = { width: bitmap.width, height: bitmap.height }
	bitmap.close()
	return result
}
async function encodeResizedImage(buffer, type, size) {
	const bitmap = await createImageBitmap(new Blob([buffer], { type }))
	try {
		let scale = Math.min(
			1,
			imageLimits.dimension / size.width,
			imageLimits.dimension / size.height,
			Math.sqrt(imageLimits.pixels / (size.width * size.height))
		)
		for (let attempt = 0; attempt < 12; attempt++) {
			const canvas = document.createElement('canvas')
			canvas.width = Math.max(1, Math.floor(size.width * scale))
			canvas.height = Math.max(1, Math.floor(size.height * scale))
			const context = canvas.getContext('2d')
			if (!context) throw new Error('Impossibile preparare l’immagine sul dispositivo.')
			context.drawImage(bitmap, 0, 0, canvas.width, canvas.height)
			const quality = type === 'image/jpeg' ? Math.max(0.5, 0.86 - attempt * 0.06) : undefined
			const output = await new Promise((resolve) => canvas.toBlob(resolve, type, quality))
			if (output?.size <= imageLimits.bytes) return output
			scale *= 0.85
		}
		throw new Error('Non riesco a ridurre questa immagine abbastanza. Scegline una diversa.')
	} finally {
		bitmap.close()
	}
}
export async function prepareImage(
	file,
	decodeDimensions = dimensions,
	encodeImage = encodeResizedImage
) {
	validateImageFile(file)
	let buffer, size
	try {
		buffer = await file.arrayBuffer()
		const bytes = new Uint8Array(buffer)
		const png = bytes[0] === 137 && bytes[1] === 80 && bytes[2] === 78 && bytes[3] === 71
		const jpeg = bytes[0] === 255 && bytes[1] === 216 && bytes[2] === 255
		if ((file.type === 'image/png' && !png) || (file.type === 'image/jpeg' && !jpeg))
			throw new Error()
		size = await decodeDimensions(buffer, file.type)
	} catch {
		throw new Error(
			'Impossibile leggere l’immagine: il file è danneggiato o non corrisponde al formato dichiarato.'
		)
	}
	if (
		!size.width ||
		!size.height ||
		size.width > imageLimits.sourceDimension ||
		size.height > imageLimits.sourceDimension ||
		size.width * size.height > imageLimits.sourcePixels
	) {
		throw new Error('L’immagine originale è troppo grande per il dispositivo.')
	}
	if (
		file.size > imageLimits.bytes ||
		size.width > imageLimits.dimension ||
		size.height > imageLimits.dimension ||
		size.width * size.height > imageLimits.pixels
	) {
		const output = await encodeImage(buffer, file.type, size)
		buffer = await output.arrayBuffer()
	}
	let binary = ''
	const bytes = new Uint8Array(buffer)
	for (let start = 0; start < bytes.length; start += 8192)
		binary += String.fromCharCode(...bytes.subarray(start, start + 8192))
	return { mimeType: file.type, name: file.name || 'Immagine', data: btoa(binary) }
}

// Deduplicate concurrent requests; immutable binary responses use the browser's private HTTP cache.
export function fetchImage(url) {
	if (!/^\/media\/[1-9]\d*$/.test(url))
		return Promise.reject(new Error('Indirizzo immagine non valido'))
	const key = `${getSessionToken()}:${url}`
	if (!pending.has(key)) {
		const request = axios.get(url, { responseType: 'blob' }).then((response) => response.data)
		pending.set(key, request)
		request
			.finally(() => {
				if (pending.get(key) === request) pending.delete(key)
			})
			.catch(() => {})
	}
	return pending.get(key)
}
