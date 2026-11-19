<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import AppAlert from './AppAlert.vue'
const props = defineProps({
	modelValue: { type: String, default: '' },
	sending: Boolean,
	attachment: Object,
	preview: String,
	error: { type: String, default: '' }
})
const emit = defineEmits(['update:modelValue', 'send', 'file', 'clear', 'dismiss-error']),
	field = ref(null),
	input = ref(null),
	form = ref(null)
function resize() {
	if (field.value) {
		field.value.style.height = 'auto'
		field.value.style.height = `${Math.min(field.value.scrollHeight, 112)}px`
	}
}
function update(event) {
	emit('update:modelValue', event.target.value)
	resize()
}
function enter(event) {
	const touchDevice = window.matchMedia?.('(pointer: coarse)').matches
	if (!event.isComposing && !event.shiftKey && !touchDevice) {
		event.preventDefault()
		emit('send')
	}
}
watch(
	() => props.modelValue,
	() =>
		nextTick(() => {
			resize()
			updateComposerHeight()
		})
)
onMounted(resize)
let viewport, orientationTimer, resizeObserver
function updateComposerHeight() {
	if (form.value)
		document.documentElement.style.setProperty('--composer-height', `${form.value.offsetHeight}px`)
}
function viewportResize() {
	const height = window.visualViewport?.height || window.innerHeight,
		viewportTop = window.visualViewport?.offsetTop || 0,
		keyboardInset = Math.max(0, window.innerHeight - height - viewportTop),
		textEntryFocused = document.activeElement?.matches(
			'input:not([type="file"]), textarea, [contenteditable="true"]'
		)
	if (!textEntryFocused || keyboardInset < 100) {
		document.documentElement.style.setProperty('--visible-height', `${window.innerHeight}px`)
		document.documentElement.style.removeProperty('--composer-bottom')
		return
	}
	document.documentElement.style.setProperty('--visible-height', `${height + viewportTop}px`)
	document.documentElement.style.setProperty('--composer-bottom', `${keyboardInset}px`)
}
function orientationChange() {
	viewportResize()
	clearTimeout(orientationTimer)
	orientationTimer = setTimeout(viewportResize, 300)
}
onMounted(() => {
	viewport = window.visualViewport
	viewportResize()
	updateComposerHeight()
	if (window.ResizeObserver) {
		resizeObserver = new ResizeObserver(updateComposerHeight)
		if (form.value) resizeObserver.observe(form.value)
	}
	viewport?.addEventListener('resize', viewportResize)
	window.addEventListener('resize', viewportResize)
	window.addEventListener('orientationchange', orientationChange)
})
onBeforeUnmount(() => {
	viewport?.removeEventListener('resize', viewportResize)
	window.removeEventListener('resize', viewportResize)
	window.removeEventListener('orientationchange', orientationChange)
	clearTimeout(orientationTimer)
	resizeObserver?.disconnect()
	document.documentElement.style.removeProperty('--visible-height')
	document.documentElement.style.removeProperty('--composer-bottom')
	document.documentElement.style.removeProperty('--composer-height')
})
</script>
<template>
	<form ref="form" class="v2-composer" @submit.prevent="emit('send')">
		<AppAlert
			v-if="error"
			class="v2-composer__error"
			:message="error"
			dismissible
			@dismiss="emit('dismiss-error')"
		/>
		<div v-if="attachment" class="attachment-preview composer-attachment">
			<img :src="preview" :alt="attachment.name" /><span>{{ attachment.name }}</span
			><button
				class="v2-button v2-button--secondary"
				type="button"
				:disabled="sending"
				@click="emit('clear')"
			>
				Rimuovi allegato
			</button>
		</div>
		<input
			ref="input"
			class="v2-sr-only"
			type="file"
			aria-label="Seleziona immagine da allegare"
			accept="image/png,image/jpeg"
			tabindex="-1"
			:disabled="sending"
			@change="emit('file', $event)"
		/>
		<button
			type="button"
			class="v2-icon-button"
			aria-label="Allega immagine"
			:disabled="sending"
			@click="input.click()"
		>
			＋
		</button>
		<label class="v2-composer__field"
			><span class="v2-sr-only">Scrivi un messaggio</span
			><textarea
				ref="field"
				:value="modelValue"
				rows="1"
				placeholder="Scrivi un messaggio…"
				:disabled="sending"
				@input="update"
				@keydown.enter="enter"
			></textarea>
		</label>
		<button
			class="v2-send-button"
			:disabled="sending || (!attachment && !modelValue.trim())"
			:aria-label="sending ? 'Invio in corso' : 'Invia messaggio'"
		>
			<svg viewBox="0 0 24 24" aria-hidden="true"><path d="m22 2-7 20-4-9-9-4zM22 2 11 13" /></svg>
		</button>
	</form>
</template>
