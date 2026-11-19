<script setup>
import { computed, ref, watch } from 'vue'
import AppAvatar from './AppAvatar.vue'
import AppDialog from './AppDialog.vue'
import ImagePicker from './ImagePicker.vue'
import AppAlert from './AppAlert.vue'
import { initials } from '../services/profile'
const props = defineProps({ conversation: Object, users: Array, busy: Boolean, error: String })
defineEmits(['close', 'rename', 'photo', 'add', 'leave', 'loadUsers'])
const title = ref(props.conversation.title),
	selected = ref('')
const available = computed(() =>
	(props.users || []).filter(
		(user) => !props.conversation.members.some((member) => member.userId === user.id)
	)
)
watch(available, (users) => {
	if (!users.some((user) => user.id === Number(selected.value))) selected.value = ''
})
</script>
<template>
	<AppDialog title="Dettagli" drawer @close="$emit('close')"
		><div class="v2-drawer__hero">
			<AppAvatar :initials="initials(conversation.title)" :src="conversation.photo" size="xlarge" />
			<h3>{{ conversation.title }}</h3>
			<p>{{ conversation.isGroup ? 'Gruppo' : 'Conversazione privata' }}</p>
		</div>
		<AppAlert :message="error" /><template v-if="conversation.isGroup"
			><section class="v2-drawer__section">
				<form class="v2-form" @submit.prevent="$emit('rename', title)">
					<label for="edit-group-name">Nome del gruppo</label
					><input
						id="edit-group-name"
						v-model.trim="title"
						maxlength="25"
						required
						:disabled="busy"
					/><button class="v2-button v2-button--primary" :disabled="busy">Salva nome</button>
				</form>
				<ImagePicker :busy="busy" label="Cambia foto gruppo" @save="$emit('photo', $event)" />
			</section>
			<section class="v2-drawer__section">
				<h3>{{ conversation.members.length }} partecipanti</h3>
				<div v-for="member in conversation.members" :key="member.userId" class="v2-member-row">
					<AppAvatar :initials="initials(member.name)" :src="member.photo" size="small" /><strong>{{
						member.name
					}}</strong>
				</div>
				<button class="v2-button v2-button--secondary" :disabled="busy" @click="$emit('loadUsers')">
					Aggiungi partecipanti
				</button>
				<form v-if="users" class="v2-form" @submit.prevent="$emit('add', Number(selected))">
					<label for="new-member">Persona da aggiungere</label
					><select id="new-member" v-model="selected" required :disabled="busy">
						<option value="" disabled>Scegli una persona</option>
						<option v-for="user in available" :key="user.id" :value="user.id">
							{{ user.name }}
						</option>
					</select>
					<p v-if="!available.length">Tutte le persone disponibili sono già nel gruppo.</p>
					<button class="v2-button v2-button--primary" :disabled="busy || !selected">
						Aggiungi
					</button>
				</form>
			</section>
			<button
				class="v2-button v2-button--danger v2-button--wide"
				:disabled="busy"
				@click="$emit('leave')"
			>
				Abbandona il gruppo
			</button></template
		>
		<section v-else class="v2-drawer__section">
			<div v-for="member in conversation.members" :key="member.userId" class="v2-member-row">
				<AppAvatar :initials="initials(member.name)" :src="member.photo" size="small" /><strong>{{
					member.name
				}}</strong>
			</div>
		</section>
	</AppDialog>
</template>
