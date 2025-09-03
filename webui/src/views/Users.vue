<script>
export default {
    data: function() {
        return {
			errormsg: null,
			ids: [],
			names: [],
			conversationId: null,
			successmsg: null,
			searchQuery: ""
		}
    },
    methods: {
        logout(){
			localStorage.removeItem("auth")
			localStorage.removeItem("id")
			this.$router.push("/login") //torno al login
		},
        settingsSection(){
			this.$router.push("/settings")
		},
        conversationsSection(){	
            this.$router.push("/conversations")
        },
		async getAllUsers() {
			this.errormsg = null	
			try{
				//console.log("Authorization Header:", localStorage.getItem("auth"));
				let response = await this.$axios.get("/users", {headers:{Authorization: localStorage.getItem("id")}});
				//console.log("Response Data: ",response.data);
				if(response) {
					this.ids = response.data.userIds;
					this.names = response.data.names;
				}
			} catch(e){
				this.errormsg = e.toString();
			}
		},
		async newConversation(userId){
			try{
				let response = await this.$axios.post(`/users/${userId}`, {}, {headers:{Authorization: localStorage.getItem("id")}});
				console.log(response.data)
				if (response) {
					this.conversationId = response.data.conversationId;

					this.successmsg = "Conversazione creata con successo!";
					setTimeout(()=> {
						this.successmsg = null;
						this.$router.push(`/conversations/${this.conversationId}`);
					},500);	
				}
			} catch(e){
				if(e.response) {
					const status = e.response.status;
					if (status === 400){
						this.errormsg = "Conversazione già esistente";
						setTimeout(()=> {
							this.errormsg = null;
						},1000);
				} else {
						this.errormsg = e.toString();
				}
				}
			}
		}
    },
	computed:{
		filteredUsers() {
			if (!this.searchQuery) {
				return this.names.map((name, index) => ({
				id: this.ids[index],
				name
				}));
			}
			const query = this.searchQuery.toLowerCase();
			return this.names
				.map((name, index) => ({ id: this.ids[index], name }))
				.filter(user => user.name.toLowerCase().includes(query));
		}
	},
    mounted() {
		this.getAllUsers();
    }
}
</script>
<template>
    <div class="navbar">
		<div class="container">
			<button type="button" class="btn" @click="conversationsSection"> Conversazioni </button>
			<button type="button" class="currentButton"> Utenti </button>
			<button type="button" class="btn" @click="settingsSection"> Impostazioni </button>
			<button type="button" class="logoutButton" @click="logout"> Logout </button>
		</div> 
	</div>
	<div class="userList">
		<div class="searchBarContainer">
			<input  class="searchBar" type="text" v-model="searchQuery" placeholder="Cerca utenti..." />
		</div>
		<!-- <UserItem v-for="(name,index) in names" 
			:key="ids[index]" :userId="ids[index]"
			:username="name"
			 @new-conversation="newConversation"
		/> -->
		<UserItem v-for="user in filteredUsers"
			:key="user.id"
			:userId="user.id"
			:username="user.name"
			@new-conversation="newConversation"
    	/>
		<SuccessMsg v-if="successmsg" :msg="successmsg" />
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>


<style scoped>
.searchBar{
	padding: 8px 12px;
	border-radius: 1rem;
	border: 2px solid rgb(4, 103, 133);
	font-family: "SF Pro";
	margin-bottom: 20px;
}
.userList{
	margin-top: 100px; /* o quanto è alta la tua navbar */
	display: flex;
	justify-content: center;
	flex-direction: column;
	align-items: center;
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
</style>