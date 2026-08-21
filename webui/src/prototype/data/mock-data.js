import {reactive} from 'vue'

export const currentUser = reactive({
	id: 1,
	name: 'Emanuele',
	initials: 'EB',
	photoUrl: ''
})

export const users = [
	{id: 2, name: 'Giulia Romano', initials: 'GR', status: 'Online', color: 'coral'},
	{id: 3, name: 'Luca Ferri', initials: 'LF', status: '', color: 'blue'},
	{id: 4, name: 'Marta Bianchi', initials: 'MB', status: 'Online', color: 'violet'},
	{id: 5, name: 'Andrea Conti', initials: 'AC', status: '', color: 'green'},
	{id: 6, name: 'Sara Moretti', initials: 'SM', status: 'Online', color: 'amber'}
]

export const conversations = [
	{
		id: 1,
		title: 'Giulia Romano',
		initials: 'GR',
		color: 'coral',
		isGroup: false,
		lastMessage: 'Perfetto, ci sentiamo più tardi!',
		lastMessageAuthor: 'Giulia',
		time: '12:42',
		unread: 2,
		online: true,
		members: [currentUser, users[0]]
	},
	{
		id: 2,
		title: 'Progetto WASA',
		initials: 'WA',
		color: 'blue',
		isGroup: true,
		lastMessage: 'Ho aggiornato la documentazione',
		lastMessageAuthor: 'Luca',
		time: '11:18',
		unread: 0,
		members: [currentUser, users[1], users[2], users[3]]
	},
	{
		id: 3,
		title: 'Calcetto del giovedì',
		initials: 'CG',
		color: 'green',
		isGroup: true,
		lastMessage: 'Campo prenotato per le 21:00 ⚽',
		lastMessageAuthor: 'Andrea',
		time: 'Ieri',
		unread: 5,
		members: [currentUser, users[1], users[3], users[4]]
	},
	{
		id: 4,
		title: 'Marta Bianchi',
		initials: 'MB',
		color: 'violet',
		isGroup: false,
		lastMessage: 'Grazie mille!',
		lastMessageAuthor: 'Tu',
		time: 'Mar',
		unread: 0,
		members: [currentUser, users[2]]
	}
]

export const messagesByConversation = {
	1: [
		{id: 1, authorId: 2, author: 'Giulia', text: 'Ciao! Come procede il progetto?', day: 'Oggi', time: '12:31', reactions: []},
		{id: 2, authorId: 1, author: 'Tu', text: 'Bene! Sto lavorando alla nuova interfaccia.', day: 'Oggi', time: '12:34', reactions: [{emoji: '✨', count: 1}]},
		{id: 3, authorId: 1, author: 'Tu', text: 'Vorrei renderla più chiara, veloce e piacevole da usare.', day: 'Oggi', time: '12:35', reactions: []},
		{id: 4, authorId: 2, author: 'Giulia', text: 'Ottima idea. La palette azzurra mi piace già molto.', day: 'Oggi', time: '12:39', reactions: [{emoji: '💙', count: 2}]},
		{id: 5, authorId: 2, author: 'Giulia', text: 'Perfetto, ci sentiamo più tardi!', day: 'Oggi', time: '12:42', reactions: []}
	],
	2: [
		{id: 1, authorId: 3, author: 'Luca', text: 'Ho preparato la prima bozza del contratto API.', day: 'Ieri', time: '18:48', reactions: []},
		{id: 2, authorId: 1, author: 'Tu', text: 'Ottimo, poi controlliamo insieme gli status code.', day: 'Ieri', time: '18:55', reactions: [{emoji: '👍', count: 3}]},
		{id: 3, authorId: 3, author: 'Luca', text: 'Ho aggiornato la documentazione', day: 'Oggi', time: '11:18', forwarded: true, reactions: []}
	],
	3: [
		{id: 1, authorId: 5, author: 'Andrea', text: 'Chi c’è questa settimana?', day: 'Ieri', time: '18:20', reactions: []},
		{id: 2, authorId: 6, author: 'Sara', text: 'Io ci sono 🙋‍♀️', day: 'Ieri', time: '18:24', reactions: [{emoji: '⚽', count: 2}]},
		{id: 3, authorId: 5, author: 'Andrea', text: 'Campo prenotato per le 21:00 ⚽', day: 'Ieri', time: '19:02', reactions: []}
	],
	4: [
		{id: 1, authorId: 4, author: 'Marta', text: 'Mi mandi il link quando hai finito?', day: 'Martedì 18 agosto', time: '16:10', reactions: []},
		{id: 2, authorId: 1, author: 'Tu', text: 'Certo, te lo mando appena è online.', day: 'Martedì 18 agosto', time: '16:12', reactions: []},
		{id: 3, authorId: 4, author: 'Marta', text: 'Grazie mille!', day: 'Martedì 18 agosto', time: '16:15', reactions: [{emoji: '😊', count: 1}]}
	]
}

export function getConversation(id) {
	return conversations.find((conversation) => conversation.id === Number(id)) || conversations[0]
}

export function getMessages(id) {
	return messagesByConversation[Number(id)] || []
}
