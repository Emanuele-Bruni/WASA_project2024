<script>
export default{
    data: function(){
        return {
            errormsg: null,
            conversationId: null,
            conversation: null,
            messages: [],
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
                console.log(response);
                this.conversation = response.data;
            }catch(e){
                this.errormsg = e.toString();
            }

        },
    },
    mounted(){
        this.getConversationInfo();
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
    <!-- Lista di messaggi (devo fare endpoint per tornare i messaggi?)-->
    <!-- DownBar per scrivere i messaggi-->
     <div class="downBar">
        <div class="downBarContent">
            <input type="text" placeholder="Scrivi un messaggio"/>
            <button class="sendButton"> Invia </button>
        </div>
     </div>
</template>

<style scoped>
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
  gap: 20px;
  border: 2px solid rgb(243, 32, 32);
  width: 75%;
  text-align: center;
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
    gap: 50px;
}
.topBar{
    background-color: white;
	padding: 13px 10px;
	position: fixed;
	left: 0;
	width: 100%;
	box-shadow: 0 5px 5px rgba(53, 81, 94, 0.5	)
}
.downBarContent{
    display: flex;
    align-items: center;
    gap: 5px;
}
.downBar{
    background-color: white;
	padding: 13px 20px;
	position: fixed;
    margin-top: 518px;
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