<script>
import EmojiPicker from '../components/EmojiPicker.vue';

export default{
    data: function(){
        return {
            errormsg: null,
            successmsg: null,
            groupError: null,

            conversationId: this.$route.params.conversationid,
            username : localStorage.getItem("username"),
            userId: localStorage.getItem("id"),
            userMap: JSON.parse(localStorage.getItem("userMap") || "{}"),
            selectedMessageId: null,
            selectedConversationId: null,
            selectedComments: [],
            specificComment: null,

            conversation: null,
            messages: [],
            newMessage: "",
            isPhoto: false,
            forwardConversations: [],
            forwardedIds: JSON.parse(localStorage.getItem("forwardedIds") || "[]"),
            membersToAdd: [],

            showChatModal: false,           //appare quando si clicca su Info
            confirmLeaveGroup: false,       //appare quando si vuole lasciare il gruppo
            confirmDeleteMessage: false,    //appare quando si vuole eliminare un messaggio
            confirmDeleteComment: false,
            showMessageOptions: false,      //appare quando si clicca sui puntini del messaggio
            forwardModal: false,            //appare quando si vuole inoltrare un messaggio
            addMemberModal: false,          //appare quando si vuole aggiungere un membro
            emojiPicker: false,             //appare quando si vuole commentare un messaggio
            allComments: false,             //appare quando si vuole eliminare un commento

        }
    },
    methods: {
        back(){
            this.$router.push("/conversations")
        },
        async getConversationInfo(){
            try{
                let response = await this.$axios.get(`/conversations/${this.conversationId}`, {headers:{Authorization: localStorage.getItem("id")}});
                //console.log(response);
                this.conversation = response.data;
            }catch(e){
                this.errormsg = e.toString();
            }
        },
        async getChatMessages(){
            try {
                let response = await this.$axios.get(`/conversations/${this.conversationId}/messages`,{headers:{Authorization: localStorage.getItem("id")}});
                // console.log(response);
                const msgs = response.data.messages;
                if (!Array.isArray(msgs) || msgs.length === 0){
                    this.errormsg = "Nessun Messaggio";
                    this.messages = []
                } else {
                    this.errormsg = null;
                    this.messages = msgs;
                }
            } catch(e){
                    this.errormsg = e.toString();
                    this.scrollToBottom();
                    this.messages = []; 
            }
        },
        async sendMessage(){
            try{
                let response = await this.$axios.post(`/conversations/${this.conversationId}/messages`,
                {text:this.newMessage},{headers:{Authorization: localStorage.getItem("id")}});
                this.newMessage = "";
                await this.getChatMessages();
                this.scrollToBottom();
                 await this.getConversationInfo();
            } catch(e) {
                if(e.response){
                    const status = e.response.status;
					if(status === 400){
						this.errormsg = "Lunghezza del messaggio non valida";
					}
                } else {
                    this.errormsg = e.toString();
                }
                this.scrollToBottom();
            }
        },
        scrollToBottom() {
            this.$nextTick(() => { 
                const el = this.$refs.bottomAnchor
                if (el) {
                    el.scrollIntoView({behavior: 'smooth'});
                }
            });
        },
        chatSettings(){
            this.showChatModal = true;
        },
        async setGroupName(editableTitle){
            try{
                let response = await this.$axios.put(`/conversations/${this.conversationId}/name`,
                {groupName: editableTitle}, {headers:{Authorization: localStorage.getItem("id")}});

                this.successmsg = "Nome gruppo aggiornato con successo";
                this.showChatModal= false;
                setTimeout(()=> {this.successmsg=null;},800);

                await this.getConversationInfo();
                this.scrollToBottom();

            } catch(e){
                if(e.response){
                    const status = e.response.status;
					if(status === 400){
						this.groupError = "Lunghezza del nome non valida";
                        setTimeout(()=> {this.groupError=null;},800);
					}
                } else {
                    this.groupError = e.toString();
                }

                this.showChatModal= false;
                this.scrollToBottom();
            }
        },
        async leaveGroup(){
            try{
                let response = await this.$axios.delete(`/conversations/${this.conversationId}/members`, {headers:{Authorization: localStorage.getItem("id")}})
                this.successmsg = "Gruppo abbandonato con successo";
                this.showChatModal = false;
                setTimeout(()=> {this.successmsg=null;},800);
                //torno alla lista delle conversazioni
                this.$router.push("/conversations")
            } catch(e) {
                this.errormsg = e.toString();
            }
        },
        openMessageOptions(messageId){
            this.selectedMessageId = messageId;
            // prendo i commenti e li salvo
            const msg = this.messages.find(m => m.messageId === messageId)
            this.selectedComments = msg ? msg.comments : []
            this.showMessageOptions = true;
        },
        async deleteMessage(messageId){
            try{
                let response = await this.$axios.delete(`/conversations/${this.conversationId}/messages/${messageId}`, {headers:{Authorization: localStorage.getItem("id")}});
                this.confirmDeleteMessage = false;
                this.showMessageOptions = false;
                await this.getChatMessages();

            } catch(e) {
                if(e.response) {
                    const status = e.response.status;
                    if(status === 403){
                        this.errormsg = "Un messaggio può essere rimosso solo dall'autore";
                        setTimeout(()=> {this.errormsg=null;},3000);
                    } else {
                        this.errormsg = e.toString();
                    }
                }
            }
        },
        async openForwardModal(messageId){
			this.selectedMessageId = messageId
			try {
				let response = await this.$axios.get("/conversations",{headers:{Authorization: localStorage.getItem("id")}});
				if(response){
					this.forwardConversations = response.data.conversations;
                    this.forwardModal = true;
				}
			} catch(e){
                this.errormsg = e.toString();
			}
		},
        async forwardMessage(convId){
            try{
                let response = await this.$axios.post(`/conversations/${this.conversationId}/messages/${this.selectedMessageId}/forward/${convId}`,
                null,{headers:{Authorization: localStorage.getItem("id")}});
                const newMessageId = response.data.messageId;
                if (!this.forwardedIds.includes(newMessageId)) {
                    this.forwardedIds.push(newMessageId);
                    // salva in localStorage
                    localStorage.setItem("forwardedIds", JSON.stringify(this.forwardedIds));
                }

                 await this.getChatMessages();

                this.successmsg = "Inoltrato con successo";
                setTimeout(()=> {this.successmsg=null;
                    this.forwardModal = false;
                    this.showMessageOptions = false;
                },1000);
            } catch(e) {
                this.errormsg = e.toString();
            }
        },
        async openAddMember(){
            try{
                let response = await this.$axios.get("/users",{headers:{Authorization: localStorage.getItem("id")}});
                if(response){
                    // Combino Ids e Names
					this.membersToAdd = response.data.userIds.map((id,index) => {
						return {id: id, name: response.data.names[index]};});
                    this.addMemberModal=true
                }
            } catch(e) {
                this.errormsg = e.toString();
            }
        },
        async addToGroup(userId){
            try{
                let response = await this.$axios.put(`/conversations/${this.conversationId}/members`,{user: userId},{headers:{Authorization: localStorage.getItem("id")}});
                this.successmsg = "Utente aggiunto con successo";
                setTimeout(()=> {this.successmsg=null;
                    this.addMemberModal = false;
                },1000);
                this.getConversationInfo();
            } catch(e) {
                this.errormsg = e.toString();
            }
        },
        async commentMessage(selectedEmoji) {
            try{
                let response = await this.$axios.post(`/conversations/${this.conversationId}/messages/${this.selectedMessageId}/comments`,
                {emoji: selectedEmoji},{headers:{Authorization: localStorage.getItem("id")}});

                // aggiorno la userMap per visualizzare subito il nome
                let storedMap = JSON.parse(localStorage.getItem("userMap"));
                storedMap[this.userId] = this.username;
                this.userMap =storedMap;

                await this.getChatMessages();
                this.emojiPicker = false;
                this.showMessageOptions = false;
                
            } catch(e) {
                this.errormsg = e.toString();
                setTimeout(()=> {this.errormsg=null;
                },1000);
            }
        },
        commentOption(commentId){
            this.confirmDeleteComment =  true;
            this.specificComment = commentId;
        },
        async uncommentMessage(commentId) {
            try{
                let response = await this.$axios.delete(`/conversations/${this.conversationId}/messages/${this.selectedMessageId}/comments/${commentId}`, {headers:{Authorization: localStorage.getItem("id")}});
                this.confirmDeleteComment = false;
                this.allComments = false;
                this.showMessageOptions=false;
                await this.getChatMessages();

            } catch(e) {
                if(e.response) {
                    this.confirmDeleteComment = false;
                    const status = e.response.status;
                    if(status === 500){
                        this.errormsg = "Un commento può essere rimosso solo dall'autore";
                        setTimeout(()=> {this.errormsg=null;},1000);
                    } else {
                        this.confirmDeleteComment = false;
                        this.errormsg = e.toString();
                    }
                }
            }
        },
        async setGroupPhoto(base64Photo){
            try{
                let response = await this.$axios.put(`/conversations/${this.conversationId}/photo`,
					{ photo: base64Photo },
					{ headers:{ Authorization: localStorage.getItem("id") } }
				);
                this.showChatModal = false;
                await this.getConversationInfo();
                this.successmsg = "Foto aggiornata con successo!";
				setTimeout(()=> {this.successmsg=null},1000);
            } catch(e){
                this.errormsg = e.toString();
            }
        },
        onFileSelected(e){
            const file = e.target.files[0];
            if(!file) return;

            const reader = new FileReader();
            reader.onload = async () => {
                const base64Photo = reader.result.split(",")[1];
                this.newMessage = "[foto]" + base64Photo;
                await this.sendMessage();
            };
            reader.readAsDataURL(file);
        },
    },
    mounted(){
        this.getConversationInfo();
        this.scrollToBottom();
        this.getChatMessages();
        this.interval = setInterval(() => this.getChatMessages(), 5000);
    },
    beforeUnmount() {
  		clearInterval(this.interval);
	}
}
</script>
<template>
    <!-- TopBar con Foto,Nome e tasto impostazioni-->
     <div class="topBar">
        <div class="topBarContent">
            <button type="button" class="backButton" @click="back"> Indietro </button>
             <div class="chatHeader" v-if="conversation">
               <img :src="'data:image/png;base64,' + conversation.photo"  class="chatPhoto"/>
               <div class="chatTitle"> {{ conversation.title }}</div>
            </div>
            <button type="button" class="chatSettingsButton" @click="chatSettings"> Info </button>
            <ChatInfoItem v-if="showChatModal"
                :photo="conversation.photo"
                :title="conversation.title"
                :isGroup="conversation.isGroup"
                :members="conversation.members"
                :conversationId="conversationId"
                @close="showChatModal = false" 
                @groupNameChanged="setGroupName"
                @left="confirmLeaveGroup=true"
                @addMember="openAddMember"
                @groupPhotoSelected="setGroupPhoto($event)"/>
            <ConfirmModal v-if="confirmLeaveGroup"
                message="Sei sicuro di voler uscire dal gruppo?"
                @yes="leaveGroup"
                @no="confirmLeaveGroup=false"
            />
            <AddToGroupModal v-if="addMemberModal"
                :conversationId="conversationId"
                :users="membersToAdd"
                @addToGroup="addToGroup($event)"
                @close="addMemberModal=false"
            />
        </div>
     </div>

    <!-- Lista di messaggi-->
    <div class="messagesSection">
        <div class="messagesList">
           <MessageItem v-for="msg in messages"
                :key="msg.messageId"
                :messageId="msg.messageId"
                :author="msg.author"
                :text="msg.text"
                :hours="msg.hours"
                :isMine="msg.authorId === parseInt(userId)"
                :comments="msg.comments"
                :isForwarded="forwardedIds.includes(msg.messageId)"
                @messageOptions="openMessageOptions"
            />
            <MessageModal v-if="showMessageOptions"
                :messageId="selectedMessageId"
                @forward="openForwardModal(selectedMessageId)"
                @delete="confirmDeleteMessage = true"
                @openEmojiPicker="emojiPicker=true"
                @openComments="allComments=true"
                @close="showMessageOptions=false"
            />
            <MessageComments v-if="allComments"
                :messageId="selectedMessageId"
                :comments="selectedComments"
                :userId="userId"
                @uncomment="commentOption($event)"
                @close="allComments=false"
            />
            <EmojiPicker v-if="emojiPicker"
                :messageId="selectedMessageId"
                @chosen="commentMessage($event)"
                @close="emojiPicker=false"
            />
            <ForwardModal v-if="forwardModal"
                :messageId="selectedMessageId"
                :conversationId="conversationId"
                :conversations="forwardConversations"
                @toConversation="forwardMessage($event)"
                @close="forwardModal=false"
            />
            <ConfirmModal v-if="confirmDeleteComment"
                message="Sei sicuro di voler eliminare il commento?"
                @yes="uncommentMessage(specificComment)"
                @no="confirmDeleteComment=false"
            />
            <ConfirmModal v-if="confirmDeleteMessage"
                message="Sei sicuro di voler eliminare il messaggio?"
                @yes="deleteMessage(selectedMessageId)"
                @no="confirmDeleteMessage=false"
            />
            <SuccessMsg v-if="successmsg" :msg="successmsg"/>
            <ErrorMsg v-if="errormsg" :msg="errormsg" :class="{ centeredError: errormsg === 'Nessun Messaggio' }"/> 
            <ErrorMsg v-if="groupError" :msg="groupError"/>
            <div ref="bottomAnchor"></div>
        </div>
    </div>

    <!-- DownBar per scrivere i messaggi-->
     <div class="downBar">
        <div class="downBarContent">
            <!-- Pulsante per scegliere la foto -->
            <button type="button" class="sendButton" @click="$refs.fileInput.click()"> 🖼️​ </button>
            <input type="file" ref="fileInput" accept="image/*" style="display: none" @change="onFileSelected" />
            <form id="sendMessageForm" @submit.prevent="sendMessage">
                <input type="text" v-model="newMessage" placeholder="Scrivi un messaggio"/>
                <button class="sendButton" type="submit"> Invia </button>
            </form>
        </div>
     </div>
</template>

<style scoped>
#sendMessageForm {
  display: flex;
  width: 100%;
  gap: 1rem;
}
.centeredError{
    text-align: center;
    width: 100%;
    margin-top: 170px;
}
.messagesSection {
    width: 100%;
    display: flex;
    justify-content: center;
    margin-top: 80px;
    margin-bottom: 80px;
}
.messagesList {
    width: 80%;
}
input{
    font-family: "SF Pro";
	width: 100%;
	padding: 0.875rem;
	border: 2px solid rgb(4, 103, 133);
	border-radius: 1rem;
	box-shadow: 1px 3px 0 rgb(4, 103, 133);
	outline: none;
	transition: ease 0.25s;
}
.sendButton{
 	font-family: "SF Pro";
	color: rgb(4, 66, 85); 
	background-color: #bad4e9;
	border: 2px solid rgb(4, 103, 133);
	border-radius: 1rem;
	padding: 0.875rem;
    width: 118px;
	box-shadow: 1px 3px 0 rgb(4, 103, 133);
}
.sendButton:hover{
    color: white;
	background-color: rgb(10, 146, 187);
}
.sendButton:active{
    box-shadow: 
	inset 0px -12px 12px -6px rgba(31, 77, 97, 0.5), 
	inset 0px 12px 12px -6px rgba(11, 51, 68, 0.5),
	inset 12px 0px 12px -6px rgba(11, 51, 68, 0.5),
	inset -12px 0px 12px -6px rgba(11, 51, 68, 0.5);
}
.chatHeader {
  display: flex;
  gap: 1rem;
  width: 80%;
}
.chatSettingsButton{
	font-family: "SF Pro";
	color: rgb(4, 66, 85); 
	background-color: #bad4e9;
	border: 2px solid rgb(4, 103, 133);
	border-radius: 1rem;
	padding: 0.875rem;
    width: 118px;
	box-shadow: 1px 3px 0 rgb(4, 103, 133);
}
.chatSettingsButton:hover{
    color: white;
	background-color: rgb(10, 146, 187);
}
.chatSettingsButton:active{
    box-shadow: 
	inset 0px -12px 12px -6px rgba(31, 77, 97, 0.5), 
	inset 0px 12px 12px -6px rgba(11, 51, 68, 0.5),
	inset 12px 0px 12px -6px rgba(11, 51, 68, 0.5),
	inset -12px 0px 12px -6px rgba(11, 51, 68, 0.5);
}
.chatPhoto {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid rgb(4, 66, 85);
  margin-left: 1rem;
}
.chatTitle {
  font-weight: bold;
  font-size: 30px;
  color: rgb(4, 66, 85);
  font-family: "SF Pro";
}
.topBarContent{
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 0.2rem;
}
.topBar{
	background-color: white;
	padding: 0.75rem 1rem;
	position: fixed;
	left: 0;
	top: 0;
	width: 100%;
	z-index: 1000;
	box-shadow: 0 5px 5px rgba(53, 81, 94, 0.5	) /* x y sfocatura colore(trasparenza)*/
}
.downBarContent{
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1.5rem;
    width: 100%;
}
.downBar{
    background-color: white;
	padding: 0.75rem 1rem;
	position: fixed;
    bottom:0;
	left: 0;
	width: 100%;
    display: flex;
	box-shadow: 0 -5px 5px rgba(53, 81, 94, 0.5	);
}
.backButton{
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
.backButton:hover{
	background-color: rgb(241, 41, 41);
}
.backButton:active{
	box-shadow: 
	inset 0px -12px 12px -6px rgba(31, 77, 97, 0.5), /* x,y,sfocatura,restringimento -> rgba(1 indica l'opacità*/
	inset 0px 12px 12px -6px rgba(11, 51, 68, 0.5),
	inset 12px 0px 12px -6px rgba(11, 51, 68, 0.5),
	inset -12px 0px 12px -6px rgba(11, 51, 68, 0.5);
}
</style>