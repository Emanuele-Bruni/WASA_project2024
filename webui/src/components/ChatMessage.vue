<script setup>
import { computed } from 'vue'
import MediaImage from './MediaImage.vue'
const props = defineProps({ message: Object, userId: Number, isGroup: Boolean })
defineEmits(['actions', 'reaction'])
const reactions = computed(() => {
	const groups = new Map()
	for (const comment of props.message.comments || []) {
		if (!groups.has(comment.emoji))
			groups.set(comment.emoji, { emoji: comment.emoji, count: 0, mine: false })
		const item = groups.get(comment.emoji)
		item.count++
		if (comment.authorId === props.userId) item.mine = true
	}
	return [...groups.values()]
})
</script>
<template>
	<article :class="['v2-message-row', { 'v2-message-row--mine': message.authorId === userId }]">
		<div class="v2-message-bubble">
			<div v-if="isGroup && message.authorId !== userId" class="v2-message-bubble__author">
				{{ message.author }}
			</div>
			<div v-if="message.forwarded" class="v2-forwarded-label">↗ Inoltrato</div>
			<MediaImage
				v-if="message.type === 'image'"
				class="v2-message-image"
				:src="message.image?.url"
				:alt="message.image?.name || 'Immagine inviata'"
				:width="message.image?.width"
				:height="message.image?.height"
			/>
			<p v-else>{{ message.text }}</p>
			<div class="v2-message-bubble__meta">
				<time :datetime="message.createdAt">{{ message.hours }}</time>
			</div>
			<button
				type="button"
				class="v2-message-menu"
				:aria-label="`Azioni sul messaggio di ${message.author}`"
				@click="$emit('actions', message.messageId)"
			>
				•••
			</button>
			<div v-if="reactions.length" class="v2-reactions">
				<button
					v-for="reaction in reactions"
					:key="reaction.emoji"
					:class="{ 'v2-reaction--mine': reaction.mine }"
					:aria-pressed="reaction.mine"
					:aria-label="`${reaction.mine ? 'Rimuovi' : 'Aggiungi'} reazione ${reaction.emoji}, ${reaction.count} persone`"
					:title="reaction.mine ? 'Rimuovi la tua reazione' : 'Aggiungi questa reazione'"
					@click="$emit('reaction', message, reaction.emoji)"
				>
					{{ reaction.emoji }} <span>{{ reaction.count }}</span>
				</button>
			</div>
		</div>
	</article>
</template>
