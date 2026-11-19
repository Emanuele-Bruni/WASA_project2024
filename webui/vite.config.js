import { fileURLToPath, URL } from 'node:url'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, process.cwd(), 'VITE_')
	return {
		plugins: [vue()],
		server: {
			host: '0.0.0.0',
			proxy: {
				'/api': {
					target: env.VITE_PROXY_TARGET || 'http://localhost:3000',
					changeOrigin: true,
					rewrite: (path) => path.replace(/^\/api/, '')
				}
			}
		},
		resolve: { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } },
		// Do not modify this constant, it is used in the evaluation.
		define: { __API_URL__: JSON.stringify(env.VITE_API_URL || 'http://localhost:3000') },
		test: { include: ['src/**/*.test.js'] }
	}
})
