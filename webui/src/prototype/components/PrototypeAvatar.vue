<script setup>
defineProps({
	initials: {type: String, required: true},
	src: {type: String, default: ''},
	color: {type: String, default: 'teal'},
	size: {type: String, default: 'medium'},
	online: {type: Boolean, default: false}
})
</script>

<template>
	<span :class="['v2-avatar', `v2-avatar--${color}`, `v2-avatar--${size}`]" aria-hidden="true">
		<img v-if="src" class="v2-avatar__image" :src="src" alt="">
		<template v-else>{{ initials }}</template>
		<span v-if="online" class="v2-avatar__presence"></span>
	</span>
</template>
