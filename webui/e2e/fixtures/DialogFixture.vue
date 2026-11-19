<script setup>
import { ref } from 'vue'
import AppDialog from '../../src/components/AppDialog.vue'
const outer = ref(false),
	inner = ref(false)
</script>
<template>
	<main class="v2-app">
		<button @click="outer = true">Apri esterno</button>
		<AppDialog v-if="outer" title="Esterno" @close="outer = false"
			><button @click="inner = true">Apri interno</button><button>Ultimo esterno</button></AppDialog
		>
		<AppDialog v-if="inner" title="Interno" @close="inner = false"
			><input aria-label="Campo interno" /><button>Ultimo interno</button></AppDialog
		>
	</main>
	<div id="dialogs" class="v2-app"></div>
</template>
