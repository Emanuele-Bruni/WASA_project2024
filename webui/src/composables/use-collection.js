import { onMounted, onBeforeUnmount, ref } from 'vue'
import { useRequestScope } from './use-request-scope'
import { createPoller } from '../services/polling'
import { fetchCollection } from '../services/server-state'
import { apiErrorMessage } from '../services/errors'
export function useCollection(name) {
	const { client, signal } = useRequestScope(),
		loading = ref(true),
		error = ref('')
	const poller = createPoller(async (pollSignal) => {
		try {
			const result = await fetchCollection(client, name, pollSignal)
			if (result) error.value = ''
		} catch (e) {
			if (!signal.aborted && e.code !== 'ERR_CANCELED') error.value = apiErrorMessage(e)
		} finally {
			if (!signal.aborted) loading.value = false
		}
	})
	onMounted(() => poller.start())
	onBeforeUnmount(() => poller.stop())
	return { client, signal, loading, error, refresh: () => poller.refresh() }
}
