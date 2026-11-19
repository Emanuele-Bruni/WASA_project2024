import js from '@eslint/js'
import vue from 'eslint-plugin-vue'
import globals from 'globals'
export default [
	{
		ignores: [
			'dist/**',
			'.yarn/**',
			'.pnp.*',
			'coverage/**',
			'playwright-report/**',
			'test-results/**'
		]
	},
	js.configs.recommended,
	...vue.configs['flat/essential'],
	// ESLint's assignment analysis does not inspect Vue templates.
	{ files: ['**/*.vue'], rules: { 'no-useless-assignment': 'off' } },
	{
		files: ['**/*.{js,vue}'],
		languageOptions: { globals: { ...globals.browser, ...globals.node, __API_URL__: 'readonly' } },
		rules: {
			'no-unused-vars': ['error', { argsIgnorePattern: '^_', caughtErrors: 'none' }],
			'vue/multi-word-component-names': 'off'
		}
	}
]
