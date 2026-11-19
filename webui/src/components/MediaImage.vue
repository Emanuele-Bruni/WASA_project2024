<script setup>
import { onBeforeUnmount, ref, watch } from 'vue'
import { fetchImage } from '../services/images'
const props = defineProps({
	src: { type: String, default: '' },
	decorative: Boolean,
	alt: { type: String, default: 'Immagine' },
	width: Number,
	height: Number
})
const objectURL = ref('')
const failed = ref(false)
let generation = 0
function release() {
	if (objectURL.value) URL.revokeObjectURL(objectURL.value)
	objectURL.value = ''
}
async function load() {
	const current = ++generation
	release()
	failed.value = false
	if (!props.src) return
	try {
		const blob = await fetchImage(props.src)
		if (current === generation) objectURL.value = URL.createObjectURL(blob)
	} catch {
		if (current === generation) failed.value = true
	}
}
watch(() => props.src, load, { immediate: true })
onBeforeUnmount(() => {
	generation++
	release()
})
</script>
<template>
	<img
		v-if="objectURL"
		:src="objectURL"
		:alt="decorative ? '' : alt"
		:width="width"
		:height="height"
		loading="lazy"
	/>
	<span
		v-else-if="failed"
		:role="decorative ? undefined : 'img'"
		:aria-hidden="decorative || undefined"
		:aria-label="alt + ': caricamento non riuscito'"
		class="media-fallback"
		>Immagine non disponibile
		<button v-if="!decorative" type="button" @click="load">Riprova</button></span
	>
	<span v-else role="img" :aria-label="alt" class="media-fallback">◌</span>
</template>
<style scoped>
.media-fallback {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	gap: 0.4rem;
	font-size: 0.75rem;
	color: #52686b;
}
</style>
