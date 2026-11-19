import { defineConfig, devices } from '@playwright/test'
export default defineConfig({
	testDir: './e2e',
	fullyParallel: true,
	forbidOnly: !!process.env.CI,
	retries: process.env.CI ? 1 : 0,
	workers: 2,
	timeout: 30000,
	expect: { timeout: 10000 },
	reporter: [['list'], ['html', { open: 'never' }]],
	use: {
		baseURL: process.env.E2E_BASE_URL || 'http://127.0.0.1:4173',
		trace: 'retain-on-failure',
		screenshot: 'only-on-failure'
	},
	projects: [
		{
			name: 'desktop-chromium',
			use: { ...devices['Desktop Chrome'], viewport: { width: 1280, height: 900 } }
		},
		{
			name: 'tablet-chromium',
			use: { ...devices['Desktop Chrome'], viewport: { width: 768, height: 1024 } }
		},
		{
			name: 'mobile-chromium',
			use: { ...devices['Pixel 7'], viewport: { width: 390, height: 844 } }
		},
		{ name: 'desktop-firefox', use: { ...devices['Desktop Firefox'] } },
		{ name: 'mobile-webkit', use: { ...devices['iPhone 13'] } }
	],
	webServer: {
		command: 'node .yarn/releases/yarn-4.5.0.cjs exec vite --host 0.0.0.0 --port 4173',
		url: 'http://127.0.0.1:4173',
		reuseExistingServer: !process.env.CI
	}
})
