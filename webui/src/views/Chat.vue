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
        mounted(){
            this.getConversationInfo();
        }
    }
}
</script>
<template>
    <!-- TopBar con Foto,Nome e tasto impostazioni-->
     <div class="topBar">
         <button type="button" class="backButton" @click="back"> Indietro </button>
         <!-- ChatItem da creare-->
          <div class="chatHeader" v-if="conversation">
            <img :src="conversation.photo" alt="Chat Photo" class="chatPhoto"/>
            <div class="chatText">
                <div class="chatTitle"> {{ conversation.title }}</div>
                <div class="chatMembers"> {{ conversation.members.name.join(", ") }}</div>

            </div>
          </div>
     </div>
    <!-- Lista di messaggi (devo fare endpoint per tornare i messaggi?)-->
    <!-- DownBar per scrivere i messaggi-->
</template>

<style scoped>
.chatHeader {
  display: flex;
  align-items: center;
  margin-left: 10px;
}

.chatPhoto {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  margin-left: 15px;
}

.chatText {
  display: flex;
  flex-direction: column;
  margin-left: 10px;
}

.chatTitle {
  font-weight: bold;
  font-size: 18px;
  color: rgb(4, 66, 85);
}

.chatMembers {
  font-size: 12px;
  color: #666;
}

.topBar{
    background-color: white;
	padding: 13px 10px;
	position: fixed;
	left: 0;
	width: 100%;
	box-shadow: 0 5px 5px rgba(53, 81, 94, 0.5	) /* x y sfocatura colore(trasparenza)*/
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