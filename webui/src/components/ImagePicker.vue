<script setup>
import { onBeforeUnmount, ref } from 'vue'
import { prepareImage } from '../services/images'
import AppAlert from './AppAlert.vue'
defineProps({ busy: Boolean, label: { type: String, default: 'Cambia immagine' } })
const emit = defineEmits(['save']),
	input = ref(null),
	image = ref(null),
	preview = ref(''),
	error = ref(''),
	reading = ref(false)
let generation = 0
function clear() {
	generation++
	if (preview.value) URL.revokeObjectURL(preview.value)
	preview.value = ''
	image.value = null
	error.value = ''
}
async function select(event) {
	const file = event.target.files?.[0]
	event.target.value = ''
	if (!file) return
	const current = ++generation
	reading.value = true
	error.value = ''
	try {
		const prepared = await prepareImage(file)
		if (current !== generation) return
		if (preview.value) URL.revokeObjectURL(preview.value)
		image.value = prepared
		preview.value = URL.createObjectURL(file)
	} catch (e) {
		if (current === generation) error.value = e.message
	} finally {
		if (current === generation) reading.value = false
	}
}
onBeforeUnmount(clear)
</script>
<template>
	<div class="image-picker">
		<input
			ref="input"
			type="file"
			aria-label="Seleziona immagine"
			accept="image/png,image/jpeg"
			class="v2-sr-only"
			tabindex="-1"
			:disabled="busy || reading"
			@change="select"
		/>
		<button
			type="button"
			class="v2-button v2-button--secondary"
			:disabled="busy || reading"
			@click="input.click()"
		>
			{{ reading ? 'Lettura…' : label }}
		</button>
		<AppAlert :message="error" />
		<div v-if="image" class="attachment-preview">
			<img :src="preview" :alt="image.name" /><span>{{ image.name }}</span>
			<button
				type="button"
				class="v2-button v2-button--primary"
				:disabled="busy"
				@click="emit('save', image)"
			>
				{{ busy ? 'Salvataggio…' : 'Salva immagine' }}
			</button>
			<button type="button" class="v2-button v2-button--secondary" :disabled="busy" @click="clear">
				Annulla immagine
			</button>
		</div>
	</div>
</template>
