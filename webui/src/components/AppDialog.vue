<script setup>
import { nextTick, onMounted, onBeforeUnmount, ref } from 'vue'
import { registerDialog } from '../services/dialog-stack'
defineProps({ title: { type: String, required: true }, drawer: Boolean })
const emit = defineEmits(['close']),
	panel = ref(null)
let dispose,
	active = true
onMounted(async () => {
	await nextTick()
	if (active) dispose = registerDialog(panel.value, () => emit('close'))
})
onBeforeUnmount(() => {
	active = false
	dispose?.()
})
</script>
<template>
	<Teleport to="#dialogs"
		><div :class="drawer ? 'v2-drawer-backdrop' : 'v2-modal-backdrop'" @click.self="emit('close')">
			<section
				ref="panel"
				:class="drawer ? 'v2-drawer' : 'v2-modal'"
				role="dialog"
				aria-modal="true"
				:aria-label="title"
				tabindex="-1"
			>
				<header class="v2-modal__header">
					<h2>{{ title }}</h2>
					<button class="v2-icon-button" type="button" aria-label="Chiudi" @click="emit('close')">
						×
					</button>
				</header>
				<slot />
			</section></div
	></Teleport>
</template>
