<script>
export default {
    data: function() {
        return {
			errormsg: null,
			username:"",
			newUsername:"",
			successmsg : null,
			userId: localStorage.getItem("id"),
			photo:"",
			selectedPhoto:"",
		}
    },
    methods: {
        logout(){
			localStorage.removeItem("auth")
			localStorage.removeItem("id")
			this.$router.push("/login") //torno al login
		},
        usersSection(){
			this.$router.push("/users")
		},
        conversationsSection(){
            this.$router.push("/conversations")
        },
		async setUsername(){
			this.successmsg = null;
			this.errormsg = null;
			const token = localStorage.getItem("id");
			if(!token){
				this.errormsg = "Utente non autenticato";
				return
			}
			try {
				let response = await this.$axios.put("/settings/username",
				{username: this.newUsername},{headers:{Authorization: localStorage.getItem("id")}});
				console.log("Response Data: ",response.data);
				if(response){
					this.successmsg = "Username aggiornato con successo!";
					localStorage.setItem("username",response.data.username);
					this.username = response.data.username;

					setTimeout(()=> {this.successmsg=null},1000);
				}
			}catch(e){
				if (e.response){
					const status = e.response.status;
					if(status===400){
						this.errormsg = "Lunghezza del nome non valida";
					} 
					if(status===409){
						this.errormsg = "Username già utilizzato";
					} else {
						this.errormsg = e.toString();
					}
				}
			}
		},
		async getMyPhoto(){
			try{
				let response = await this.$axios.get(`/settings/photo/${this.userId}`,{headers:{Authorization: localStorage.getItem("id")}});
				console.log(response.data);
				this.photo = response.data.photo;
				localStorage.setItem("myPhoto",this.photo);
			} catch(e) {
				this.errormsg = e.toString();
			}
		},
		onFileChange(e){
			const file = e.target.files[0];
			if(!file) return;
			const reader = new FileReader();
			reader.onload = () => {
				this.selectedPhoto = reader.result.split(",")[1];
				this.setMyPhoto();
				this.photo = this.selectedPhoto;
			};
			reader.readAsDataURL(file);
		},
		async setMyPhoto(){
			try {
				let response = await this.$axios.put("/settings/photo",
					{ photo: this.selectedPhoto },
					{ headers:{ Authorization: localStorage.getItem("id") } }
				);
				this.successmsg = "Foto aggiornata con successo!";
				localStorage.setItem("myPhoto", this.photo);
				setTimeout(()=> {this.successmsg=null},1000);
			} catch(e) {
				this.errormsg = e.toString();
			}
		}
    },
    mounted() {
		const savedUsername = localStorage.getItem("username") || "";
		this.username = savedUsername;
		this.newUsername = savedUsername;
		this.getMyPhoto();
    }
}
</script>
<template>
    <div class="navbar">
		<div class="container">
			<button type="button" class="btn" @click="conversationsSection"> Conversazioni </button>
			<button type="button" class="btn" @click="usersSection"> Utenti </button>
			<button type="button" class="currentButton"> Impostazioni </button>
			<button type="button" class="logoutButton" @click="logout"> Logout </button>
		</div> 
	</div>
	<div class="subNavbar">
		<div class="settings">
			<p>Modifica la tua foto </p>
			<!-- CERCHIO PER ANTEPRIMA DELLA FOTO -->
			<div class="avatar-wrapper">
				<img :src="`data:image/png;base64,` + photo" class="avatar" />
			</div>
			<!-- SELEZIONE DELLA NUOVA FOTO -->
			<input type="file" accept="image/*" @change="onFileChange">
	
			<p>Modifica il tuo username </p>
			<SuccessMsg v-if="successmsg" :msg="successmsg" />
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
			<!-- BARRA PER IMPOSTARE USERNAME-->
			<form id="newUsernameForm" @submit.prevent="setUsername">
				<input type="text" id="username" v-model="newUsername" :placeholder="username"/>
				<button class="saveButton" type="submit"> Salva </button>
			</form>	<br>
		</div>
	</div>

</template>


<style scoped>
.subNavbar{
	width: 100%;
	display: flex;
	justify-content: center;
}
.avatar-wrapper {
	display: flex;
	justify-content: center;
  	align-items: center;
  	margin-bottom: 1rem;
}
.avatar {
  	width: 200px;
  	height: 200px;
  	border-radius: 50%;
  	object-fit: cover;
	border: 3px solid rgb(4, 103, 133);
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
input[type="file"] {
	font-family: "SF Pro";
	width: 70%;              /* stessa larghezza del bottone */
	max-width: 300px;
	padding: 0.6rem;
	border: 2px solid rgb(4, 103, 133);
	border-radius: 1rem;
	box-shadow: 1px 3px 0 rgb(4, 103, 133);
	text-align: center;
	cursor: pointer;
}
p{
	font-size: 1.4rem;
	font-family: "SF Pro";
	color: rgb(4, 78, 100); 
	text-align: center;
	margin: 0.5rem 0;
}
.settings{
	width: 50%;
	border: 1px solid black;
	margin-top: 80px;
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 0.5rem;
	gap: 1rem;
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
.saveButton {
	text-align: center;
	font-family: "SF Pro";
	color: rgb(4, 66, 85); 
	background-color: #bad4e9;
	border: 2px solid rgb(4, 103, 133);
	border-radius: 1rem;
	padding: 0.875rem;
	width: 100%;
	margin-top: 0.3rem;
	box-shadow: 1px 3px 0 rgb(4, 103, 133);
}
.saveButton:hover{
	background-color:rgb(10, 146, 187);
	color: white;
}
.saveButton:active{
	box-shadow: 
	inset 0px -12px 12px -6px rgba(31, 77, 97, 0.5), /* x,y,sfocatura,restringimento -> rgba(1 indica l'opacità*/
	inset 0px 12px 12px -6px rgba(11, 51, 68, 0.5),
	inset 12px 0px 12px -6px rgba(11, 51, 68, 0.5),
	inset -12px 0px 12px -6px rgba(11, 51, 68, 0.5);
}
</style>