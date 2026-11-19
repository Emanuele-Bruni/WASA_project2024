<script setup>
defineProps({ message: String, success: Boolean, retry: Boolean, dismissible: Boolean })
defineEmits(['retry', 'dismiss'])
</script>
<template>
	<div
		v-if="message"
		:class="['v2-alert', success ? 'v2-alert--success' : 'v2-alert--error']"
		:role="success ? 'status' : 'alert'"
	>
		<span>{{ message }}</span
		><button
			v-if="retry"
			class="v2-button v2-button--secondary"
			type="button"
			@click="$emit('retry')"
		>
			Riprova
		</button>
		<button
			v-if="dismissible"
			class="v2-alert__close"
			aria-label="Chiudi avviso"
			type="button"
			@click="$emit('dismiss')"
		>
			×
		</button>
	</div>
</template>
