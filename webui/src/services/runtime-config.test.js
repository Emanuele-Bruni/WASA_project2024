import { afterEach, expect, it, vi } from 'vitest'

afterEach(() => {
	vi.unstubAllGlobals()
	vi.resetModules()
})

it('uses the API URL supplied when the container starts', async () => {
	vi.resetModules()
	vi.stubGlobal('__WASA_CONFIG__', { apiURL: '/api' })
	const { default: axios } = await import('./axios')
	expect(axios.defaults.baseURL).toBe('/api')
})

it('uses the Vite API URL when runtime settings are absent', async () => {
	vi.resetModules()
	vi.stubGlobal('__WASA_CONFIG__', {})
	const { default: axios } = await import('./axios')
	expect(axios.defaults.baseURL).toBe('http://localhost:3000')
})
