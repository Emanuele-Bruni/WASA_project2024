<script>
export default {
	data: function() {
		return {
			errormsg: null,
			successmsg: null,
			conversations: [],
			userMap: JSON.parse(localStorage.getItem("userMap") || "{}"),
			userId: localStorage.getItem("id"),
			lastMessages: JSON.parse(localStorage.getItem("lastMessages") || "{}"),
			showGroupModal: false,
			groupTitle: "",
			userOptions: [],
			selectedUserIds: [],
			groupError: null,
			groupId: null,
		}
	},
	methods: {
		logout(){
			localStorage.removeItem("auth")
			localStorage.removeItem("id")
			this.$router.push("/login") //torno al login
		},
		usersSection(){
			//vado in un'altra pagina
			this.$router.push("/users") //path della prossima pagina
		},
		settingsSection(){
			this.$router.push("/settings")
		},
		async getMyConversations(){
			this.errormsg = null
			try {
				let response = await this.$axios.get("/conversations",{headers:{Authorization: localStorage.getItem("id")}});
				console.log(response.data);
				if(response){
					this.conversations = response.data.conversations;
				}
			} catch(e){
				if(e.response) {
					const status = e.response.status;
					if (status === 404){
						this.errormsg = "Nessuna Conversazione";
				} else {
						this.errormsg = e.toString();
				}
				}
			}
		},
		chatSection(conversationId){
			this.$router.push(`/conversations/${conversationId}`)
		},
		async newGroup(){
			this.groupTitle="";
			this.selectedUserIds=[];
			this.groupError = null;
			try{
				let response = await this.$axios.get("/users",{headers:{Authorization: localStorage.getItem("id")}});
				//console.log(response.data);
				if(response && response.data.userIds && response.data.names) {
					// Combino Ids e Names
					this.userOptions = response.data.userIds.map((id,index) => {
						return {id: id, name: response.data.names[index]};});
					this.showGroupModal = true;
				}
				await this.getMyConversations();
			}catch(e) {
				this.groupError = e.toString();
			}
		},
		async createGroup(){
			this.groupError = null;
			if (this.groupTitle.trim() === "" || this.selectedUserIds.length < 2) {
				this.groupError = "Inserisci un titolo valido e seleziona almeno due utenti.";
				return;
			}
			try {
				let response = await this.$axios.post("/conversations",
					{title: this.groupTitle.trim(), users: this.selectedUserIds,},
					{headers:{Authorization: localStorage.getItem("id")}});
				this.successmsg = "Gruppo creato con successo!";
				console.log(response.data);
				this.groupId = response.data.groupId;
				setTimeout(()=> {
						this.successmsg = null;
						this.showGroupModal = false;
						this.$router.push(`/conversations/${this.groupId}`);
					},700);
				 await this.getMyConversations();
			} catch(e){
				if (e.response){
					const status = e.response.status;
					if(status === 400){
						this.groupError = "Nome gruppo troppo lungo";
					}
				} else {
					this.groupError = e.toString();
				}
				await this.getMyConversations();
			}
		},
		getLastMessage(conversationId){
			return this.lastMessages[conversationId] || null;
		},
		async refreshLastMessages(){
			try {
				for(const conv of this.conversations){
					let response = await this.$axios.get(`/conversations/${conv.conversationId}/messages`, 
														{ headers: { Authorization: localStorage.getItem("id") }});
					const msgs = response.data.messages;
					if(Array.isArray(msgs) && msgs.length){
						const lastMsg = msgs[msgs.length - 1];
						let displayText = lastMsg.text.startsWith("[foto]") ? "📷Foto" : lastMsg.text;
						if(displayText.length > 10) displayText = displayText.substring(0,10) + "...";
						this.lastMessages[conv.conversationId] = {
							author: this.userMap[lastMsg.authorId] || lastMsg.author,
							text: displayText,
							hours: lastMsg.hours
						};
					}
				}
				localStorage.setItem("lastMessages", JSON.stringify(this.lastMessages));
			} catch(e){
				console.error("Errore aggiornando ultimi messaggi", e);
			}
		}

	},
	mounted() {
		this.getMyConversations();
		this.interval = setInterval(() => {this.getMyConversations(); this.refreshLastMessages();}, 30000);
	},
	beforeUnmount() {
  		clearInterval(this.interval);
	}
}
</script>

<template>
	<div>
		<div class="navbar">
			<div class="container">
				<button type="button" class="currentButton"> Conversazioni </button>
				<button type="button" class="btn" @click="usersSection"> Utenti </button>
				<button type="button" class="btn" @click="settingsSection"> Impostazioni </button>
				<button type="button" class="logoutButton" @click="logout"> Logout </button>
			</div> 
		</div>
		
		<div class="below-navbar">
			<!-- Pulsante per creare nuovo gruppo -->
			<button type="button" class="btnNewGroup" @click="newGroup"> Nuovo Gruppo </button>
		</div>

		<!-- Lista di tutte le conversazioni dell'utente -->
		 <div class="listWrapper">
			 <div class="conversations-list">
				 <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg> 
				 <div>
					 <ConversationItem v-for="conv in conversations" 
						 :key="conv.conversationId" 
						 :photo="conv.photo" 
						 :title="conv.title" 
						 :isGroup = "conv.isGroup"
						 :lastMessage ="getLastMessage(conv.conversationId)"
						 @click="chatSection(conv.conversationId)"></ConversationItem>
				 </div>
			 </div>
		 </div>

		 <!-- Card che appare una volta premuto il pulsante per il nuovo gruppo -->
		 <div v-if="showGroupModal" class="modal-overlay">
			<div class="modal-content">
				<h2> Nuovo Gruppo </h2>
				<form @submit.prevent="createGroup">
					<label> Titolo: </label> 
					<input type="text" v-model="groupTitle" placeholder="Titolo del gruppo"/>

					<label> Seleziona utenti: </label>
					<div class="user-list">
						<div v-for="user in userOptions" :key="user.id" class="user-option">
							<input
								type="checkbox"
								:id="user.id"
								:value="user.id"
								v-model="selectedUserIds"
							/>
							<label :for="user.id">{{ user.name }}</label>
						</div>
					</div>

					<div class="modal-buttons">
						<button> Crea </button>
						<button @click="showGroupModal = false"> Annulla </button>
					</div>

					<ErrorMsg v-if="groupError" :msg="groupError"></ErrorMsg>
					<SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>
				</form>
			</div>
		 </div>

	</div>
</template>

<style scoped>
.listWrapper{
	width: 100%;
	display: flex;
	justify-content: center;
}
.conversations-list{
	width: 80%;
}
.below-navbar{
	margin-top: 100px; /* o quanto è alta la tua navbar */
	display: flex;
	justify-content: center;
	position: sticky;
}
.btnNewGroup{
	text-align: center;
	font-family: "SF Pro";
	color: rgb(4, 66, 85); 
	background-color: #bad4e9;
	border: 2px solid rgb(4, 103, 133);
	border-radius: 1rem;
	padding: 0.5rem;
	width: 80%;
	box-shadow: 1px 3px 0 rgb(4, 103, 133);
	position: sticky;
}
.btnNewGroup:hover{
	color: white;
	background-color: rgb(10, 146, 187);
}
.btnNewGroup:active{
	box-shadow: 
	inset 0px -12px 12px -6px rgba(31, 77, 97, 0.5), /* x,y,sfocatura,restringimento -> rgba(1 indica l'opacità*/
	inset 0px 12px 12px -6px rgba(11, 51, 68, 0.5),
	inset 12px 0px 12px -6px rgba(11, 51, 68, 0.5),
	inset -12px 0px 12px -6px rgba(11, 51, 68, 0.5);
}	
.navbar{
	background-color: rgb(178, 207, 233);
	padding: 0.75rem 1rem;
	position: fixed;
	left: 0;
	top: 0;
	width: 100%;
	z-index: 1000;
	box-shadow: 0 5px 5px rgba(53, 81, 94, 0.5	) /* x y sfocatura colore(trasparenza)*/
}
.container{
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 0.2rem;
}
.logoutButton{	
	color: white;
	background-color: rgb(235, 85, 85);
	text-align: center;
	font-family: "SF Pro";
	padding: 0.875rem;
	border: 2px solid rgb(243, 32, 32);
	border-radius: 1rem;
	width: 118px;
	box-shadow: 1px 3px 0 rgb(243, 32, 32);
}
.logoutButton:hover{
	background-color: rgb(241, 41, 41);
}
.logoutButton:active{
	box-shadow: 
	inset 0px -12px 12px -6px rgba(31, 77, 97, 0.5), /* x,y,sfocatura,restringimento -> rgba(1 indica l'opacità*/
	inset 0px 12px 12px -6px rgba(11, 51, 68, 0.5),
	inset 12px 0px 12px -6px rgba(11, 51, 68, 0.5),
	inset -12px 0px 12px -6px rgba(11, 51, 68, 0.5);
}
.btn{
	text-align: center;
	font-family: "SF Pro";
	color: rgb(4, 66, 85); 
	background-color: #bad4e9;
	border: 2px solid rgb(4, 103, 133);
	border-radius: 1rem;
	padding: 0.875rem;
	width: 118px;
	box-shadow: 1px 3px 0 rgb(4, 103, 133);
}
.btn:hover{
	color: white;
	background-color: rgb(10, 146, 187);
}
.btn:active{
	box-shadow: 
	inset 0px -12px 12px -6px rgba(31, 77, 97, 0.5), /* x,y,sfocatura,restringimento -> rgba(1 indica l'opacità*/
	inset 0px 12px 12px -6px rgba(11, 51, 68, 0.5),
	inset 12px 0px 12px -6px rgba(11, 51, 68, 0.5),
	inset -12px 0px 12px -6px rgba(11, 51, 68, 0.5);
}
.currentButton{
	text-align: center;
	font-family: "SF Pro";
	color: rgb(4, 66, 85); 
	background-color: rgb(199, 177, 177);
	border: 2px solid rgb(4, 103, 133);
	border-radius: 1rem;
	padding: 0.875rem;
	width: 118px;
	box-shadow: 1px 3px 0 rgb(4, 103, 133);
}
.modal-overlay {
	position: fixed;
	top: 0; left: 0;
	width: 100%;
	height: 100%;
	background: rgba(0, 0, 0, 0.6);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 9999;
}

.modal-content {
	background: white;
	border-radius: 20px;
	padding: 20px;
	width: 90%;
	max-width: 400px;
	box-shadow: 0 5px 20px rgba(0, 0, 0, 0.3);
	font-family: "SF Pro", sans-serif;
}

.modal-content h2 {
	margin-bottom: 10px;
	font-family: "SF Pro";
  	color: rgb(4, 103, 133);
}

.modal-content label {
	display: block;
	margin-top: 10px;
	font-weight: bold;
  	color: rgb(4, 103, 133);
}

.modal-content input[type="text"] {
	padding: 8px;
	margin-top: 3px;
	border-radius: 5px;
	border: 1px solid #ccc;
}

.user-list {
	max-height: 200px;
	overflow-y: auto;
	margin-top: 5px;
}

.user-option {
	margin-bottom: 5px;
	gap: 8px;
	align-items: center;
	display: flex;
}


.modal-buttons {
	display: flex;
	justify-content: space-between;
	margin-top: 20px;
}

.modal-buttons button {
	padding: 10px 20px;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	font-family: "SF Pro";
}

.modal-buttons button:first-child {
	background-color: rgb(4, 103, 133);
	color: white;
}

.modal-buttons button:last-child {
	background-color: rgb(226, 22, 22);
	color: white;
}


</style>
