<script>
import MessageItem from '../components/MessageItem.vue';

export default{
    data: function(){
        return {
            errormsg: null,
            successmsg: null,
            conversationId: null,
            conversation: null,
            messages: [],
            newMessage: "",
            username : localStorage.getItem("username"),
            userId: localStorage.getItem("id"),
        }
    },
    methods: {
        back(){
            this.$router.push("/conversations")
        },
        async getConversationInfo(){
            this.conversationId = this.$route.params.conversationid;
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
                console.log(response);
                const msgs = response.data.Messages;
                if (!Array.isArray(msgs) || msgs.length === 0){
                    this.errormsg = "Nessun Messaggio";
                    this.messages = []
                } else {
                    this.errormsg = null;
                    this.messages = msgs;
                }
                console.log(this.errormsg);
                this.scrollToBottom();
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
            } catch(e) {
                this.errormsg = e.toString();
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
        }
    },
    mounted(){
        this.getConversationInfo();
        this.getChatMessages();
    }
}
</script>
<template>
    <!-- TopBar con Foto,Nome e tasto impostazioni-->
     <div class="topBar">
        <div class="topBarContent">
            <button type="button" class="backButton" @click="back"> Indietro </button>
             <div class="chatHeader" v-if="conversation">
               <img :src="conversation.photo"  class="chatPhoto"/>
               <div class="chatTitle"> {{ conversation.title }}</div>
            </div>
            <button type="button" class="chatSettingsButton"> Impostazioni</button>
        </div>
     </div>

    <!-- Lista di messaggi (devo fare endpoint per tornare i messaggi)-->
    <div class="messagesSection">
        <div class="messagesList">
           <MessageItem v-for="msg in messages"
                :key="msg.messageId"
                :messageId="msg.messageId"
                :author="msg.author"
                :text="msg.text"
                :hours="msg.hours"
                :isMine="msg.authorId === parseInt(userId)"
            />
            <ErrorMsg v-if="errormsg" :msg="errormsg" :class="{ centeredError: errormsg === 'Nessun Messaggio' }"/> 
            <div ref="bottomAnchor"></div>
        </div>
    </div>

    <!-- DownBar per scrivere i messaggi-->
     <div class="downBar">
        <div class="downBarContent">
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
  border: 2px solid rgb(243, 32, 32);
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
    border: 2px solid black;
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