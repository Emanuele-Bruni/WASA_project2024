import { onBeforeUnmount } from 'vue'
import axios from '../services/axios'
import { scopedClient } from '../services/request-scope'
export function useRequestScope() {
	const controller = new AbortController()
	onBeforeUnmount(() => controller.abort())
	return { client: scopedClient(axios, controller.signal), signal: controller.signal }
}
