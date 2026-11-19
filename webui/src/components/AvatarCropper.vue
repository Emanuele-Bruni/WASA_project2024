<script setup>
import { onBeforeUnmount, ref, watch } from 'vue'
import { fetchImage, prepareImage } from '../services/images'
import AppAlert from './AppAlert.vue'

const props = defineProps({ src: String, busy: Boolean })
const emit = defineEmits(['save', 'cancel'])
const preview = ref(''),
	error = ref(''),
	loading = ref(true),
	saving = ref(false)
const x = ref(50),
	y = ref(50)
let bitmap,
	generation = 0

async function load() {
	const current = ++generation
	loading.value = true
	error.value = ''
	try {
		const blob = await fetchImage(props.src)
		const nextBitmap = await createImageBitmap(blob)
		if (current !== generation) return nextBitmap.close()
		bitmap?.close()
		bitmap = nextBitmap
		if (preview.value) URL.revokeObjectURL(preview.value)
		preview.value = URL.createObjectURL(blob)
	} catch {
		if (current === generation) error.value = 'Non riesco a caricare la foto profilo.'
	} finally {
		if (current === generation) loading.value = false
	}
}
watch(() => props.src, load, { immediate: true })
onBeforeUnmount(() => {
	generation++
	bitmap?.close()
	if (preview.value) URL.revokeObjectURL(preview.value)
})

async function saveCrop() {
	if (!bitmap || saving.value || props.busy) return
	saving.value = true
	error.value = ''
	try {
		const side = Math.min(bitmap.width, bitmap.height)
		const sx = (bitmap.width - side) * (x.value / 100)
		const sy = (bitmap.height - side) * (y.value / 100)
		const canvas = document.createElement('canvas')
		canvas.width = 1024
		canvas.height = 1024
		const context = canvas.getContext('2d')
		if (!context) throw new Error('Impossibile ritagliare la foto sul dispositivo.')
		context.drawImage(bitmap, sx, sy, side, side, 0, 0, canvas.width, canvas.height)
		const blob = await new Promise((resolve) => canvas.toBlob(resolve, 'image/jpeg', 0.88))
		if (!blob) throw new Error('Impossibile preparare il ritaglio.')
		const image = await prepareImage(new File([blob], 'foto-profilo.jpg', { type: 'image/jpeg' }))
		emit('save', image)
	} catch (e) {
		error.value = e.message || 'Non riesco a preparare il ritaglio.'
	} finally {
		saving.value = false
	}
}
</script>

<template>
	<section class="avatar-cropper" aria-label="Ritaglia foto profilo">
		<p>Sposta la foto per scegliere cosa mostrare nell’avatar.</p>
		<AppAlert :message="error" />
		<div v-if="loading" role="status">Caricamento foto…</div>
		<div v-else-if="preview" class="avatar-cropper__preview">
			<img :src="preview" alt="Anteprima del ritaglio" :style="{ objectPosition: `${x}% ${y}%` }" />
		</div>
		<label>Orizzontale <input v-model.number="x" type="range" min="0" max="100" /></label>
		<label>Verticale <input v-model.number="y" type="range" min="0" max="100" /></label>
		<div class="v2-form__actions">
			<button
				class="v2-button v2-button--primary"
				:disabled="loading || saving || busy"
				@click="saveCrop"
			>
				{{ saving || busy ? 'Salvataggio…' : 'Salva ritaglio' }}
			</button>
			<button
				class="v2-button v2-button--secondary"
				:disabled="saving || busy"
				@click="emit('cancel')"
			>
				Annulla
			</button>
		</div>
	</section>
</template>
